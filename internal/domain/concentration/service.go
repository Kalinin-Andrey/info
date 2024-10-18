package concentration

import (
	"context"
	"info/internal/domain"
	"time"

	validation "github.com/go-ozzo/ozzo-validation/v4"
)

type CmcApi interface {
	GetAnalytics(ctx context.Context, CurrencyID uint, Range string) (*ConcentrationList, error)
}

type Service struct {
	replicaSet ReplicaSet
	cmcApi     CmcApi
}

func NewService(replicaSet ReplicaSet, cmcApi CmcApi) *Service {
	return &Service{
		replicaSet: replicaSet,
		cmcApi:     cmcApi,
	}
}

const (
	defaultCapacity = 100

	TimeRange_1M  = "1M"
	TimeRange_1Y  = "1Y"
	TimeRange_All = "All"
)

var TimeRangeList = []interface{}{
	TimeRange_1M,
	TimeRange_1Y,
	TimeRange_All,
}

func TimeRangeValidate(s string) error {
	return validation.Validate(s, validation.Required, validation.In(TimeRangeList...))
}

func (s *Service) Upsert(ctx context.Context, entity *Concentration) error {
	return s.replicaSet.WriteRepo().Upsert(ctx, entity)
}

func (s *Service) ImportTx(ctx context.Context, tx domain.Tx, currencyID uint, importLastTime *time.Time) (err error) {
	if importLastTime == nil || time.Now().Add(-time.Hour*24*365).After(*importLastTime) {
		if err = s.importTx(ctx, tx, currencyID, TimeRange_All); err != nil {
			return err
		}
	}

	if importLastTime == nil || time.Now().Add(-time.Hour*24*31).After(*importLastTime) {
		if err = s.importTx(ctx, tx, currencyID, TimeRange_1Y); err != nil {
			return err
		}
	}

	if err = s.importTx(ctx, tx, currencyID, TimeRange_1M); err != nil {
		return err
	}

	return nil
}

func (s *Service) importTx(ctx context.Context, tx domain.Tx, currencyID uint, timeRange string) (err error) {
	if err = TimeRangeValidate(timeRange); err != nil {
		return err
	}

	item, err := s.cmcApi.GetAnalytics(ctx, currencyID, timeRange)
	if err != nil {
		return err
	}

	if err = s.replicaSet.WriteRepo().MUpsertTx(ctx, tx, item.Slice()); err != nil {
		return err
	}
	return nil
}
