package oracul_analytics

import (
	"context"
	"info/internal/domain/currency"
	"info/internal/domain/oracul_daily_balance_stats"
	"info/internal/domain/oracul_holder_stats"
	"info/internal/domain/oracul_speedometers"
	"time"
)

type OraculAnalyticsAPIClient interface {
	GetHoldersStats(ctx context.Context, currencyID uint, blockchain string, coinAddress string) (*ImportData, error)
}

type ImportData struct {
	OraculAnalytics             *OraculAnalytics
	OraculSpeedometers          *oracul_speedometers.OraculSpeedometers
	OraculHolderStats           *oracul_holder_stats.OraculHolderStats
	OraculDailyBalanceStatsList *oracul_daily_balance_stats.OraculDailyBalanceStatsList
}

type Service struct {
	replicaSet               ReplicaSet
	oraculAnalyticsAPIClient OraculAnalyticsAPIClient
	oraculSpeedometers       *oracul_speedometers.Service
	oraculHolderStats        *oracul_holder_stats.Service
	oraculDailyBalanceStats  *oracul_daily_balance_stats.Service
}

func NewService(replicaSet ReplicaSet, oraculAnalyticsAPIClient OraculAnalyticsAPIClient, oraculSpeedometers *oracul_speedometers.Service, oraculHolderStats *oracul_holder_stats.Service, oraculDailyBalanceStats *oracul_daily_balance_stats.Service) *Service {
	return &Service{
		replicaSet:               replicaSet,
		oraculAnalyticsAPIClient: oraculAnalyticsAPIClient,
		oraculSpeedometers:       oraculSpeedometers,
		oraculHolderStats:        oraculHolderStats,
		oraculDailyBalanceStats:  oraculDailyBalanceStats,
	}
}

func (s *Service) Create(ctx context.Context, entity *OraculAnalytics) error {
	return s.replicaSet.WriteRepo().Upsert(ctx, entity)
}

func (s *Service) Import(ctx context.Context, tokenAddressList *currency.TokenAddressList) (err error) {
	if tokenAddressList == nil || len(*tokenAddressList) == 0 {
		return nil
	}
	var tokenAddress currency.TokenAddress
	var importData *ImportData

	for _, tokenAddress = range *tokenAddressList {
		importData, err = s.oraculAnalyticsAPIClient.GetHoldersStats(ctx, tokenAddress.CurrencyID, tokenAddress.Blockchain, tokenAddress.Address)
		if err != nil {
			return err
		}
		if err = s.upsertImportData(ctx, importData); err != nil {
			return err
		}
		time.Sleep(5 * time.Second)
	}

	return nil
}

func (s *Service) upsertImportData(ctx context.Context, importData *ImportData) (err error) {

	if importData.OraculAnalytics != nil {
		if err = s.replicaSet.WriteRepo().Upsert(ctx, importData.OraculAnalytics); err != nil {
			return err
		}
	}

	if importData.OraculSpeedometers != nil {
		if err = s.oraculSpeedometers.Create(ctx, importData.OraculSpeedometers); err != nil {
			return err
		}
	}

	if importData.OraculHolderStats != nil {
		if err = s.oraculHolderStats.Create(ctx, importData.OraculHolderStats); err != nil {
			return err
		}
	}

	if importData.OraculDailyBalanceStatsList != nil && len(*importData.OraculDailyBalanceStatsList) > 0 {
		if err = s.oraculDailyBalanceStats.MCreate(ctx, importData.OraculDailyBalanceStatsList); err != nil {
			return err
		}
	}

	return nil
}
