package currency

import "context"

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
