package currency

import (
	"context"
)

type ReplicaSet interface {
	WriteRepo() WriteRepository
	ReadRepo() ReadRepository
}

type WriteRepository interface {
	Create(ctx context.Context, entity *Currency) (ID uint, err error)
	Update(ctx context.Context, entity *Currency) error
	Delete(ctx context.Context, ID uint) error
}

type ReadRepository interface {
	Get(ctx context.Context, currencyID uint) (*Currency, error)
	MGet(ctx context.Context, currencyIDs *[]uint) (*[]Currency, error)
	GetAll(ctx context.Context) (*[]Currency, error)
}
