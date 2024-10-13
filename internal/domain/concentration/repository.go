package concentration

import (
	"context"
)

type ReplicaSet interface {
	WriteRepo() WriteRepository
	ReadRepo() ReadRepository
}

type WriteRepository interface {
	Upsert(ctx context.Context, entity *Concentration) error
	MUpsert(ctx context.Context, entities *[]Concentration) error
}

type ReadRepository interface {
	Get(ctx context.Context, currencyID uint) (*Concentration, error)
	MGet(ctx context.Context, currencyIDs *[]uint) (*[]Concentration, error)
	GetAll(ctx context.Context) (*[]Concentration, error)
}
