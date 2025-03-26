package app

import (
	"context"
	"errors"
	"info/internal/domain/concentration"
	"info/internal/domain/portfolio_item"
	"info/internal/domain/price_and_cap"
	"info/internal/integration"
	"info/internal/pkg/config"
	"log"

	"info/internal/domain/currency"
	"info/internal/infrastructure"
	"info/internal/infrastructure/repository/tsdb_cluster"
)

type App struct {
	config      *config.AppConfig
	Infra       *infrastructure.Infrastructure
	Integration *integration.Integration
	Domain      *Domain
}

type Domain struct {
	Currency      *currency.Service
	PriceAndCap   *price_and_cap.Service
	Concentration *concentration.Service
	PortfolioItem *portfolio_item.Service
}

// New func is a constructor for the App
func New(ctx context.Context, cfg *config.Configuration) *App {
	log.Println("Core app is starting...")
	log.Println("infrastructure start create...")
	infr, err := infrastructure.New(ctx, cfg.App.InfraAppConfig(), cfg.Infra)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("done")

	log.Println("integration start create...")
	integr, err := integration.New(&integration.AppConfig{
		NameSpace:   cfg.App.NameSpace,
		Subsystem:   cfg.App.Name,
		Service:     cfg.App.Service,
		Environment: cfg.App.Environment,
	}, cfg.Integration, infr.Logger)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("done")

	app := &App{
		config:      cfg.App,
		Infra:       infr,
		Integration: integr,
	}

	app.SetupServices()

	return app
}

func (app *App) SetupServices() {
	app.Domain = &Domain{
		PriceAndCap:   price_and_cap.NewService(tsdb_cluster.NewPriceAndCapReplicaSet(app.Infra.TsDB), app.Integration.CmcApi),
		Concentration: concentration.NewService(tsdb_cluster.NewConcentrationReplicaSet(app.Infra.TsDB), app.Integration.CmcApi),
		PortfolioItem: portfolio_item.NewService(tsdb_cluster.NewPortfolioItemReplicaSet(app.Infra.TsDB), app.Integration.CmcApi),
	}
	app.Domain.Currency = currency.NewService(tsdb_cluster.NewCurrencyReplicaSet(app.Infra.TsDB), app.Domain.PriceAndCap, app.Domain.Concentration, app.Integration.CmcApi, app.Integration.CmcProApi)
}

func (app *App) Run() error {
	return nil
}

func (app *App) Stop() error {
	return errors.Join(
		app.Integration.Close(),
		app.Infra.Close(),
	)
}
