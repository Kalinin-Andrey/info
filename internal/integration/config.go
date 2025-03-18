package integration

import (
	"info/internal/integration/cmc_api"
	"info/internal/integration/cmc_pro_api"
)

type Config struct {
	CmcApi    *cmc_api.Config
	CmcProApi *cmc_pro_api.Config
}

type UsageConfig struct {
	CmcApi    bool
	CmcProApi bool
}
