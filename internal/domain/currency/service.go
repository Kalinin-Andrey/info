package currency

import (
	"context"
	"errors"
	"fmt"
	"info/internal/domain"
	"info/internal/pkg/apperror"
	"strconv"
)

type Service struct {
	replicaSet ReplicaSet
}

func NewService(replicaSet ReplicaSet) *Service {
	return &Service{
		replicaSet: replicaSet,
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

func (s *Service) Import(ctx context.Context, listOfCurrencySlugs []string) error {
	var tx domain.Tx
	var err error
	var rv uint
	var list *NmDimensionsList

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

	if tx, rv, err = s.cluster.GetNmDimensionsMaxRvTx(ctx); err != nil && !errors.Is(err, apperror.ErrNotFound) {
		return err
	}

	if list, rv, err = s.apiClient.GetNmDimensionsList(ctx, rv); err != nil {
		if errors.Is(err, apperror.ErrNotFound) {
			return nil
		}
		return err
	}
	wblogger.Info(ctx, "apiClient.GetNmDimensionsList: "+strconv.Itoa(len(*list)))

	if err = s.cluster.MUpsertNmDimensions(ctx, list.Slice()); err != nil {
		return err
	}

	if err = s.cluster.UpdateNmDimensionsMaxRvTx(ctx, tx, rv); err != nil {
		return err
	}

	return nil
}

func (s *Service) baseImport(ctx context.Context, listOfCurrencySlugs []string) error {
	
}
