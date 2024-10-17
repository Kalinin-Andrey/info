package concentration

import (
	"context"
)

type CmcApi interface {
	GetAnalytics(ctx context.Context, CurrencyID uint, Range string) (*ConcentrationList, error)
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

func (s *Service) Upsert(ctx context.Context, entity *Concentration) error {
	return s.replicaSet.WriteRepo().Upsert(ctx, entity)
}

func (s *Service) Get(ctx context.Context, currencyID uint) (*Concentration, error) {
	return s.replicaSet.ReadRepo().Get(ctx, currencyID)
}

func (s *Service) GetAll(ctx context.Context) (*[]Concentration, error) {
	return s.replicaSet.ReadRepo().GetAll(ctx)
}
