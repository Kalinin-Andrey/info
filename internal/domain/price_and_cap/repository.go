package price_and_cap

import (
	"context"
)

type ReplicaSet interface {
	WriteRepo() WriteRepository
	ReadRepo() ReadRepository
}

type WriteRepository interface {
	Upsert(ctx context.Context, entity *PriceAndCap) (err error)
	MUpsert(ctx context.Context, entities *[]PriceAndCap) error
}

type ReadRepository interface {
	MGet(ctx context.Context, currencyIDs *[]uint) (PriceAndCapMap, error)
}
