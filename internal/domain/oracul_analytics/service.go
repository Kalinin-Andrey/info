package oracul_analytics

import (
	"context"
	"info/internal/domain/oracul_daily_balance_stats"
	"info/internal/domain/oracul_holder_stats"
	"info/internal/domain/oracul_speedometers"
)

type OraculAnalyticsAPIClient interface {
	Import(ctx context.Context) (oraculAnalytics *OraculAnalytics, oraculSpeedometers *oracul_speedometers.OraculSpeedometers, oraculHolderStats *oracul_holder_stats.OraculHolderStats, oraculDailyBalanceStatsList *oracul_daily_balance_stats.OraculDailyBalanceStatsList)
}

type ImportData struct {
	OraculAnalytics             *OraculAnalytics
	OraculSpeedometers          *oracul_speedometers.OraculSpeedometers
	OraculHolderStats           *oracul_holder_stats.OraculHolderStats
	OraculDailyBalanceStatsList *oracul_daily_balance_stats.OraculDailyBalanceStatsList
}

type Service struct {
	replicaSet              ReplicaSet
	oraculSpeedometers      *oracul_speedometers.Service
	oraculHolderStats       *oracul_holder_stats.Service
	oraculDailyBalanceStats *oracul_daily_balance_stats.Service
}

func NewService(replicaSet ReplicaSet, oraculSpeedometers *oracul_speedometers.Service, oraculHolderStats *oracul_holder_stats.Service, oraculDailyBalanceStats *oracul_daily_balance_stats.Service) *Service {
	return &Service{
		replicaSet:              replicaSet,
		oraculSpeedometers:      oraculSpeedometers,
		oraculHolderStats:       oraculHolderStats,
		oraculDailyBalanceStats: oraculDailyBalanceStats,
	}
}

func (s *Service) Create(ctx context.Context, entity *OraculAnalytics) error {
	return s.replicaSet.WriteRepo().Upsert(ctx, entity)
}
