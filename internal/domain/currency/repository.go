package currency

import (
	"context"
	"info/internal/domain"
)

type ReplicaSet interface {
	WriteRepo() WriteRepository
	ReadRepo() ReadRepository
}

type WriteRepository interface {
	Begin(ctx context.Context) (domain.Tx, error)
	GetImportMaxTimeForUpdate(ctx context.Context, tx domain.Tx, currencyIDs *[]uint) (*[]ImportMaxTime, error)
	Create(ctx context.Context, entity *Currency) (ID uint, err error)
	Update(ctx context.Context, entity *Currency) error
	Delete(ctx context.Context, ID uint) error
}

type ReadRepository interface {
	Get(ctx context.Context, ID uint) (*Currency, error)
	GetBySlug(ctx context.Context, slug string) (*Currency, error)
	GetImportMaxTime(ctx context.Context, currencyIDs *[]uint) (*[]ImportMaxTime, error)
	MGet(ctx context.Context, IDs *[]uint) (*CurrencyList, error)
	MGetBySlug(ctx context.Context, slugs *[]string) (*CurrencyList, error)
	GetAll(ctx context.Context) (*[]Currency, error)
}
