package config

import (
	"flag"
	"fmt"
	"info/internal/integration"
	"time"

	"github.com/spf13/viper"

	"info/internal/infrastructure"
)

const (
	defaultPathToConfig   = "/etc/config/app.yml"
	defaultPathToMSConfig = "/etc/config/ms.yml"

	Environment_Local = "local"
	Environment_Dev   = "dev"
	Environment_Stage = "stage"
	Environment_Prd   = "prd"
)

type AppConfig struct {
	NameSpace   string
	Name        string
	Service     string
	Environment string
}

func (c *AppConfig) InfraAppConfig() *infrastructure.AppConfig {
	return &infrastructure.AppConfig{
		NameSpace:   c.NameSpace,
		Name:        c.Name,
		Service:     c.Service,
		Environment: c.Environment,
	}
}

type Configuration struct {
	App         *AppConfig
	API         *API
	Cli         *CliConfig
	Integration *integration.Config
	Infra       *infrastructure.Config
}

type API struct {
	Rest    *RestAPIConfig
	Metrics *RestAPIConfig
	Probes  *RestAPIConfig
}

type RestAPIConfig struct {
	Addr         string
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
	IdleTimeout  time.Duration
}

type CliConfig struct {
	CurrencyCollector *CurrencyCollector
}

type CurrencyCollector struct {
	Duration            time.Duration
	PortfolioSourceIDs  []string
	ListOfCurrencySlugs []string
}

// Get func return the app config
func Get() (*Configuration, error) {
	// config is the app config
	var config Configuration = Configuration{}
	// pathToConfig is a path to the app config
	var pathToConfig string
	var pathToMSConfig string

	viper.AutomaticEnv() // read in environment variables that match
	//viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	//viper.BindEnv("pathToConfig")
	defPathToConfig := defaultPathToConfig
	if viper.Get("pathToConfig") != nil {
		defPathToConfig = viper.Get("pathToConfig").(string)
	}

	flag.StringVar(&pathToConfig, "config", defPathToConfig, "path to YAML/JSON config file")
	flag.Parse()

	if err := config.readConfig(pathToConfig); err != nil {
		return &config, err
	}

	defPathToMSConfig := defaultPathToMSConfig
	if viper.Get("pathToMSConfig") != nil {
		defPathToMSConfig = viper.Get("pathToMSConfig").(string)
	}
	flag.StringVar(&pathToMSConfig, "msconfig", defPathToMSConfig, "path to YAML/JSON config file for micro-service")
	flag.Parse()

	if err := config.readMSConfig(pathToMSConfig); err != nil {
		return &config, err
	}

	// Всё, что ниже в этой функции, должно автоматом подхватываться, но почему-то на кубере не сработало. Пока придётся в рукопашную.
	var ok bool
	var st, varName string
	var val interface{}

	if config.Integration.CmcApi != nil {
		varName = "INTEGRATION_CMCAPI_COOKIE"
		val = viper.Get(varName)
		if st, ok = val.(string); !ok {
			return nil, fmt.Errorf("Can not convert "+varName+" = %v to string.", val)
		} else {
			config.Integration.CmcApi.Cookie = st
		}
	}

	if config.Integration.CmcProApi != nil {
		varName = "INTEGRATION_CMCPROAPI_TOKEN"
		val = viper.Get(varName)
		if st, ok = val.(string); !ok {
			return nil, fmt.Errorf("Can not convert "+varName+" = %v to string.", val)
		} else {
			config.Integration.CmcProApi.Token = st
		}
	}

	return &config, nil
}

func (c *Configuration) readConfig(pathToConfig string) error {
	viper.SetConfigFile(pathToConfig)

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return fmt.Errorf("Config file not found in %q", pathToConfig)
		} else {
			return fmt.Errorf("Config file was found in %q, but was produced error: %w", pathToConfig, err)
		}
	}

	err := viper.Unmarshal(c)
	if err != nil {
		return fmt.Errorf("Config unmarshal error: %w", err)
	}
	return nil
}
