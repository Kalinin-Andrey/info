package concentration

import (
	"context"
)

type ReplicaSet interface {
	WriteRepo() WriteRepository
	ReadRepo() ReadRepository
}

type WriteRepository interface {
	Upsert(ctx context.Context, entity *Concentration) (err error)
	MUpsert(ctx context.Context, entities *[]Concentration) error
}

type ReadRepository interface {
	MGet(ctx context.Context, currencyIDs *[]uint) (ConcentrationMap, error)
}
