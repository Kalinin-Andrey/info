package oracul_analytics_api

import (
	"info/internal/domain/oracul_analytics"
	"info/internal/domain/oracul_daily_balance_stats"
	"info/internal/domain/oracul_holder_stats"
	"info/internal/domain/oracul_speedometers"
	"time"
)

type GetHoldersStatsResponse struct {
	WhalesConcentration float64               `json:"whales_concentration"`
	WormIndex           float64               `json:"worm_index"`
	GrowthFuel          float64               `json:"growth_fuel"`
	Speedometers        *Speedometers         `json:"speedometers"`
	HolderStats         *HolderStats          `json:"holder_stats"`
	DailyBalanceStats   DailyBalanceStatsList `json:"daily_balance_stats"`
}

func (r *GetHoldersStatsResponse) ImportData(currencyID uint, ts time.Time) (*oracul_analytics.ImportData, error) {
	oraculDailyBalanceStatsList, err := r.DailyBalanceStats.OraculDailyBalanceStatsList(currencyID)
	if err != nil {
		return nil, err
	}
	return &oracul_analytics.ImportData{
		OraculAnalytics:             r.OraculAnalytics(currencyID, ts),
		OraculSpeedometers:          r.Speedometers.OraculSpeedometers(currencyID, ts),
		OraculHolderStats:           r.HolderStats.OraculHolderStats(currencyID, ts),
		OraculDailyBalanceStatsList: oraculDailyBalanceStatsList,
	}, nil
}

func (r *GetHoldersStatsResponse) OraculAnalytics(currencyID uint, ts time.Time) *oracul_analytics.OraculAnalytics {
	return &oracul_analytics.OraculAnalytics{
		CurrencyID:          currencyID,
		WhalesConcentration: r.WhalesConcentration,
		WormIndex:           r.WormIndex,
		GrowthFuel:          r.GrowthFuel,
		Ts:                  ts,
	}
}

type Speedometers struct {
	Whales    *SpeedometersItem `json:"whales"`
	Investors *SpeedometersItem `json:"investors"`
	Retailers *SpeedometersItem `json:"retailers"`
}
type SpeedometersItem struct {
	BuyRate  float64 `json:"buy_rate"`
	SellRate float64 `json:"sell_rate"`
	Volume   float64 `json:"volume"`
}

func (e *Speedometers) OraculSpeedometers(currencyID uint, ts time.Time) *oracul_speedometers.OraculSpeedometers {
	return &oracul_speedometers.OraculSpeedometers{
		CurrencyID:        currencyID,
		WhalesBuyRate:     e.Whales.BuyRate,
		WhalesSellRate:    e.Whales.SellRate,
		WhalesVolume:      e.Whales.Volume,
		InvestorsBuyRate:  e.Investors.BuyRate,
		InvestorsSellRate: e.Investors.SellRate,
		InvestorsVolume:   e.Investors.Volume,
		RetailersBuyRate:  e.Retailers.BuyRate,
		RetailersSellRate: e.Retailers.SellRate,
		RetailersVolume:   e.Retailers.Volume,
		Ts:                ts,
	}
}

type HolderStats struct {
	Whales    *HolderStatsItem `json:"whales"`
	Investors *HolderStatsItem `json:"investors"`
	Retailers *HolderStatsItem `json:"retailers"`
}
type HolderStatsItem struct {
	Volume       float64 `json:"volume"`
	TotalHolders uint    `json:"total_holders"`
}

func (e *HolderStats) OraculHolderStats(currencyID uint, ts time.Time) *oracul_holder_stats.OraculHolderStats {
	return &oracul_holder_stats.OraculHolderStats{
		CurrencyID:            currencyID,
		WhalesVolume:          e.Whales.Volume,
		WhalesTotalHolders:    e.Whales.TotalHolders,
		InvestorsVolume:       e.Investors.Volume,
		InvestorsTotalHolders: e.Investors.TotalHolders,
		RetailersVolume:       e.Retailers.Volume,
		RetailersTotalHolders: e.Retailers.TotalHolders,
		Ts:                    ts,
	}
}

type DailyBalanceStatsList map[string]DailyBalanceStats
type DailyBalanceStats struct {
	Whales    *DailyBalanceStatsItem `json:"whales"`
	Investors *DailyBalanceStatsItem `json:"investors"`
	Retailers *DailyBalanceStatsItem `json:"retailers"`
}
type DailyBalanceStatsItem struct {
	Balance      float64 `json:"balance"`
	TotalHolders uint    `json:"total_holders"`
}

func (e *DailyBalanceStatsList) OraculDailyBalanceStatsList(currencyID uint) (*oracul_daily_balance_stats.OraculDailyBalanceStatsList, error) {
	if e == nil || len(*e) == 0 {
		return nil, nil
	}
	var err error
	var date string
	var d time.Time
	var item DailyBalanceStats
	var resItem *oracul_daily_balance_stats.OraculDailyBalanceStats
	res := make(oracul_daily_balance_stats.OraculDailyBalanceStatsList, 0, len(*e))

	for date, item = range *e {
		d, err = time.Parse(time.DateOnly, date)
		if err != nil {
			return nil, err
		}
		resItem = item.OraculDailyBalanceStats(currencyID, d)
		res = append(res, *resItem)
	}

	return &res, nil
}

func (e *DailyBalanceStats) OraculDailyBalanceStats(currencyID uint, d time.Time) *oracul_daily_balance_stats.OraculDailyBalanceStats {
	return &oracul_daily_balance_stats.OraculDailyBalanceStats{
		CurrencyID:            currencyID,
		WhalesBalance:         e.Whales.Balance,
		WhalesTotalHolders:    e.Whales.TotalHolders,
		InvestorsBalance:      e.Investors.Balance,
		InvestorsTotalHolders: e.Investors.TotalHolders,
		RetailersBalance:      e.Retailers.Balance,
		RetailersTotalHolders: e.Retailers.TotalHolders,
		D:                     d,
	}
}
