package concentration

import (
	"context"
)

type ReplicaSet interface {
	WriteRepo() WriteRepository
	ReadRepo() ReadRepository
}

type WriteRepository interface {
	Create(ctx context.Context, entity *Concentration) (ID uint, err error)
	Update(ctx context.Context, entity *Concentration) error
	Delete(ctx context.Context, CurrencyID uint) error
}

type ReadRepository interface {
	Get(ctx context.Context, currencyID uint) (*Concentration, error)
	MGet(ctx context.Context, currencyIDs *[]uint) (*[]Concentration, error)
	GetAll(ctx context.Context) (*[]Concentration, error)
}
