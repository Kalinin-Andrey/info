package currency

import (
	"context"
	"errors"
	"fmt"
	"info/internal/domain"
	"info/internal/domain/concentration"
	"info/internal/domain/price_and_cap"
	"info/internal/pkg/apperror"
)

type CmcApi interface {
	GetCurrency(ctx context.Context, currencySlug string) (*Currency, error)
}

type Service struct {
	replicaSet    ReplicaSet
	priceAndCap   *price_and_cap.Service
	concentration *concentration.Service
	cmcApi        CmcApi
}

func NewService(replicaSet ReplicaSet, priceAndCap *price_and_cap.Service, concentration *concentration.Service, cmcApi CmcApi) *Service {
	return &Service{
		replicaSet:    replicaSet,
		priceAndCap:   priceAndCap,
		concentration: concentration,
		cmcApi:        cmcApi,
	}
}

func (s *Service) Create(ctx context.Context, entity *Currency) (ID uint, err error) {
	return s.replicaSet.WriteRepo().Create(ctx, entity)
}

func (s *Service) Update(ctx context.Context, entity *Currency) error {
	return s.replicaSet.WriteRepo().Update(ctx, entity)
}

func (s *Service) Delete(ctx context.Context, ID uint) error {
	return s.replicaSet.WriteRepo().Delete(ctx, ID)
}

func (s *Service) Get(ctx context.Context, ID uint) (*Currency, error) {
	return s.replicaSet.ReadRepo().Get(ctx, ID)
}

func (s *Service) GetAll(ctx context.Context) (*[]Currency, error) {
	return s.replicaSet.ReadRepo().GetAll(ctx)
}

func (s *Service) Import(ctx context.Context, listOfCurrencySlugs *[]string) (err error) {
	var tx domain.Tx

	currencyList, err := s.baseImport(ctx, listOfCurrencySlugs)
	if err != nil {
		return err
	}

	tx, err = s.replicaSet.WriteRepo().Begin(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("[%w] Recover from panic: %v", apperror.ErrInternal, r)
		}

		if err == nil {
			if err = tx.Commit(ctx); err == nil {
				return
			}
			err = fmt.Errorf("[%w] Commit error: %w", apperror.ErrInternal, err)
		}

		if tx != nil {
			if err2 := tx.Rollback(ctx); err2 != nil {
				err = errors.Join(err, fmt.Errorf("[%w] Rollback error: %w", apperror.ErrInternal, err))
			}
		}

	}()

	importMaxTimeList, err := s.replicaSet.WriteRepo().GetImportMaxTimeForUpdate(ctx, tx, currencyList.IDs())
	if err != nil {
		return err
	}

	var importMaxTimeItem ImportMaxTime
	for _, importMaxTimeItem = range *importMaxTimeList {
		if err = s.priceAndCap.ImportTx(ctx, tx, importMaxTimeItem.CurrencyID, importMaxTimeItem.PriceAndCap); err != nil {
			return err
		}
		if err = s.concentration.ImportTx(ctx, tx, importMaxTimeItem.CurrencyID, importMaxTimeItem.Concentration); err != nil {
			return err
		}
	}

	return nil
}

func (s *Service) baseImport(ctx context.Context, listOfCurrencySlugs *[]string) (currencyList *CurrencyList, err error) {
	if listOfCurrencySlugs == nil || len(*listOfCurrencySlugs) == 0 {
		return nil, nil
	}

	exists, err := s.replicaSet.ReadRepo().MGetBySlug(ctx, listOfCurrencySlugs)
	if err != nil {
		return nil, err
	}
	notExistsSlugs := make([]string, 0, len(*listOfCurrencySlugs))
	existsSlugsMap := make(map[string]struct{})
	var item Currency
	for _, item = range *exists {
		existsSlugsMap[item.Slug] = struct{}{}
	}

	var ok bool
	var slug string
	for _, slug = range *listOfCurrencySlugs {
		if _, ok = existsSlugsMap[slug]; !ok {
			notExistsSlugs = append(notExistsSlugs, slug)
		}
	}

	if len(notExistsSlugs) > 0 {
		importedCurrencyList, err := s.baseImportBySlug(ctx, &notExistsSlugs)
		if err != nil {
			return nil, err
		}
		(*exists) = append(*exists, (*importedCurrencyList)...)
	}
	return exists, nil
}

func (s *Service) baseImportBySlug(ctx context.Context, listOfCurrencySlugs *[]string) (importedCurrencyList *CurrencyList, err error) {
	if listOfCurrencySlugs == nil || len(*listOfCurrencySlugs) == 0 {
		return nil, nil
	}

	var slug string
	var item *Currency
	importedCurrency := make(CurrencyList, len(*listOfCurrencySlugs))
	for _, slug = range *listOfCurrencySlugs {
		if item, err = s.cmcApi.GetCurrency(ctx, slug); err != nil {
			return nil, err
		}
		_, err = s.replicaSet.WriteRepo().Create(ctx, item)
		if err != nil {
			return nil, err
		}
		importedCurrency = append(importedCurrency, *item)
	}
	return &importedCurrency, nil
}
