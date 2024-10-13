package price_and_cap

import "context"

type Service struct {
	replicaSet ReplicaSet
}

func NewService(replicaSet ReplicaSet) *Service {
	return &Service{
		replicaSet: replicaSet,
	}
}

func (s *Service) Create(ctx context.Context, entity *PriceAndCap) (ID uint, err error) {
	return s.replicaSet.WriteRepo().Upsert(ctx, entity)
}

func (s *Service) Update(ctx context.Context, entity *PriceAndCap) error {
	return s.replicaSet.WriteRepo().Update(ctx, entity)
}

func (s *Service) Delete(ctx context.Context, currencyID uint) error {
	return s.replicaSet.WriteRepo().Delete(ctx, currencyID)
}

func (s *Service) Get(ctx context.Context, currencyID uint) (*PriceAndCap, error) {
	return s.replicaSet.ReadRepo().Get(ctx, currencyID)
}

func (s *Service) GetAll(ctx context.Context) (*[]PriceAndCap, error) {
	return s.replicaSet.ReadRepo().GetAll(ctx)
}
