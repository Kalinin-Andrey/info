package oracul_analytics_api

import (
	"info/internal/domain/oracul_analytics"
	"info/internal/domain/oracul_daily_balance_stats"
	"info/internal/domain/oracul_holder_stats"
	"info/internal/domain/oracul_speedometers"
)

type GetHoldersStatsResponse struct {
	WhalesConcentration float64               `json:"whales_concentration"`
	WormIndex           float64               `json:"worm_index"`
	GrowthFuel          float64               `json:"growth_fuel"`
	Speedometers        *Speedometers         `json:"speedometers"`
	HolderStats         *HolderStats          `json:"holder_stats"`
	DailyBalanceStats   DailyBalanceStatsList `json:"daily_balance_stats"`
}

func (r *GetHoldersStatsResponse) ImportData() *oracul_analytics.ImportData {

}

func (r *GetHoldersStatsResponse) OraculAnalytics() *oracul_analytics.OraculAnalytics {

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

func (e *Speedometers) OraculSpeedometers() *oracul_speedometers.OraculSpeedometers {

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

func (e *HolderStats) OraculHolderStats() *oracul_holder_stats.OraculHolderStats {

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

func (e *DailyBalanceStatsList) OraculDailyBalanceStatsList() *oracul_daily_balance_stats.OraculDailyBalanceStatsList {

}

func (e *DailyBalanceStats) OraculDailyBalanceStats() *oracul_daily_balance_stats.OraculDailyBalanceStats {

}
