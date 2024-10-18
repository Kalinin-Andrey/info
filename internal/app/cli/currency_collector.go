package cli

import (
	"go.uber.org/zap"
	"info/internal/pkg/config"
	"math/rand"

	"context"
	"time"

	"github.com/spf13/cobra"
)

// currencyCollector ...
var currencyCollector = &cobra.Command{
	Use:   "currency-collector",
	Short: "It is the currency-collector command.",
	Long:  `It is the currency-collector command: consumer for collecting currencies.`,
	Run: func(cmd *cobra.Command, args []string) {
		CliApp.currencyCollector(cmd, args)
	},
}

func (app *App) currencyCollector(cmd *cobra.Command, args []string) {
	cfg := app.config.CurrencyCollector

	go func() {
		time.Sleep(time.Duration(rand.Int63n(int64(time.Now().Second()%30)+1)) * time.Second)
		app.dimensionCollector_Exec(app.ctx, cfg)
		ticker := time.NewTicker(cfg.Duration)
		defer ticker.Stop()

		for {
			select {
			case <-app.ctx.Done():
				return
			case <-ticker.C:
				app.dimensionCollector_Exec(app.ctx, cfg)
			}
		}
	}()

	return
}

func (app *App) dimensionCollector_Exec(ctx context.Context, cfg *config.CurrencyCollector) {
	app.Infra.Logger.Info("PriceAndCap.Import: starts iteration...")

	if err := app.Domain.Currency.Import(ctx, &cfg.ListOfCurrencySlugs); err != nil {
		app.Infra.Logger.Info("PriceAndCap.Import: iteration completed with errors!", zap.Error(err))
		return
	}
	app.Infra.Logger.Info("PriceAndCap.Import: iteration completed successfully!")
}
