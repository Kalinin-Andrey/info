package integration

import (
	"errors"
	"go.uber.org/zap"
	"info/internal/integration/cmc_api"
	"info/internal/integration/cmc_pro_api"
)

type AppConfig struct {
	NameSpace   string
	Subsystem   string
	Service     string
	Environment string
}

type Integration struct {
	CmcApi    *cmc_api.CmcApiClient
	CmcProApi *cmc_pro_api.CmcApiClient
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
