package price_and_cap

import (
	"context"
)

type ReplicaSet interface {
	WriteRepo() WriteRepository
	ReadRepo() ReadRepository
}

type WriteRepository interface {
	Create(ctx context.Context, entity *PriceAndCap) (ID uint, err error)
	Update(ctx context.Context, entity *PriceAndCap) error
	Delete(ctx context.Context, CurrencyID uint) error
}

type ReadRepository interface {
	Get(ctx context.Context, currencyID uint) (*PriceAndCap, error)
	MGet(ctx context.Context, currencyIDs *[]uint) (*[]PriceAndCap, error)
	GetAll(ctx context.Context) (*[]PriceAndCap, error)
}
