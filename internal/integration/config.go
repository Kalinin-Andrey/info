package integration

import "info/internal/integration/cmc_api"

type Config struct {
	CmcApi *cmc_api.Config
}

type UsageConfig struct {
	CmcApi bool
}
