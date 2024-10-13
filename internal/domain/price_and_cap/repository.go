package price_and_cap

import (
	"context"
)

type ReplicaSet interface {
	WriteRepo() WriteRepository
	ReadRepo() ReadRepository
}

type WriteRepository interface {
	Upsert(ctx context.Context, entity *PriceAndCap) error
	MUpsert(ctx context.Context, entities *[]PriceAndCap) error
}

type ReadRepository interface {
	Get(ctx context.Context, currencyID uint) (*PriceAndCap, error)
	MGet(ctx context.Context, currencyIDs *[]uint) (*[]PriceAndCap, error)
	GetAll(ctx context.Context) (*[]PriceAndCap, error)
}
