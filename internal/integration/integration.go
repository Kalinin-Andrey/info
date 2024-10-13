package integration

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"info/internal/domain/concentration"
	"info/internal/domain/price_and_cap"
	"info/internal/integration/cmc_api"
)

type AppConfig struct {
	NameSpace   string
	Subsystem   string
	Service     string
	Environment string
}

type CmcApi interface {
	GetDetailChart(ctx context.Context, CurrencyID uint, Range string) (*price_and_cap.PriceAndCapList, error)
	GetAnalytics(ctx context.Context, CurrencyID uint, Range string) (*concentration.ConcentrationList, error)
}

type Integration struct {
	CmcApi CmcApi
}

func New(appConfig *AppConfig, cfg *Config, logger *zap.Logger) (*Integration, error) {
	integration := &Integration{}

	if cfg.CmcApi != nil {
		integration.CmcApi = cmc_api.New(&cmc_api.AppConfig{
			NameSpace: appConfig.NameSpace,
			Subsystem: appConfig.Subsystem,
			Service:   appConfig.Service,
		}, cfg.CmcApi, logger)
	}

	return integration, nil
}

func (intgr *Integration) Close() error {
	return errors.Join()
}
