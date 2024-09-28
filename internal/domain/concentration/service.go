package concentration

import "context"

type Service struct {
	replicaSet ReplicaSet
}

func NewService(replicaSet ReplicaSet) *Service {
	return &Service{
		replicaSet: replicaSet,
	}
}

func (s *Service) Create(ctx context.Context, entity *Concentration) (ID uint, err error) {
	return s.replicaSet.WriteRepo().Create(ctx, entity)
}

func (s *Service) Update(ctx context.Context, entity *Concentration) error {
	return s.replicaSet.WriteRepo().Update(ctx, entity)
}

func (s *Service) Delete(ctx context.Context, currencyID uint) error {
	return s.replicaSet.WriteRepo().Delete(ctx, currencyID)
}

func (s *Service) Get(ctx context.Context, currencyID uint) (*Concentration, error) {
	return s.replicaSet.ReadRepo().Get(ctx, currencyID)
}

func (s *Service) GetAll(ctx context.Context) (*[]Concentration, error) {
	return s.replicaSet.ReadRepo().GetAll(ctx)
}
