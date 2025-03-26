package integration

import (
	"context"
	"errors"
	"go.uber.org/zap"
	"info/internal/domain/concentration"
	"info/internal/domain/currency"
	"info/internal/domain/portfolio_item"
	"info/internal/domain/price_and_cap"
	"info/internal/integration/cmc_api"
	"info/internal/integration/cmc_pro_api"
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
	GetCurrency(ctx context.Context, currencySlug string) (*currency.Currency, error)
	GetPortfolioSummary(ctx context.Context, portfolioSourceId string) (*portfolio_item.PortfolioItemList, error)
}

type CmcProApi interface {
	GetCurrenciesByIDs(ctx context.Context, currencyIDs *[]uint) (currencyMap currency.CurrencyMap, err error)
	GetCurrenciesBySlugs(ctx context.Context, slugs *[]string) (currencyMap currency.CurrencyMap, err error)
}

type Integration struct {
	CmcApi    CmcApi
	CmcProApi CmcProApi
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

	if cfg.CmcProApi != nil {
		integration.CmcProApi = cmc_pro_api.New(&cmc_pro_api.AppConfig{
			NameSpace: appConfig.NameSpace,
			Subsystem: appConfig.Subsystem,
			Service:   appConfig.Service,
		}, cfg.CmcProApi, logger)
	}

	return integration, nil
}

func (intgr *Integration) Close() error {
	return errors.Join()
}
