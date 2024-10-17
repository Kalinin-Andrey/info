package price_and_cap

import (
	"context"
)

type CmcApi interface {
	GetDetailChart(ctx context.Context, CurrencyID uint, Range string) (*PriceAndCapList, error)
}

type Service struct {
	replicaSet ReplicaSet
	cmcApi     CmcApi
}

func NewService(replicaSet ReplicaSet, cmcApi CmcApi) *Service {
	return &Service{
		replicaSet: replicaSet,
		cmcApi:     cmcApi,
	}
}

func (s *Service) Upsert(ctx context.Context, entity *PriceAndCap) error {
	return s.replicaSet.WriteRepo().Upsert(ctx, entity)
}

func (s *Service) Get(ctx context.Context, currencyID uint) (*PriceAndCap, error) {
	return s.replicaSet.ReadRepo().Get(ctx, currencyID)
}

func (s *Service) GetAll(ctx context.Context) (*[]PriceAndCap, error) {
	return s.replicaSet.ReadRepo().GetAll(ctx)
}
