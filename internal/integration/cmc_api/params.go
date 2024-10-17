package cmc_api

import (
	"info/internal/domain/concentration"
	"info/internal/domain/currency"
	"info/internal/domain/price_and_cap"
	"info/internal/pkg/apperror"
	"strconv"
	"time"
)

type Status struct {
	Timestamp    string `json:"timestamp"`
	ErrorCode    string `json:"error_code"`
	ErrorMessage string `json:"error_message"`
	Elapsed      string `json:"elapsed"`
	CreditCount  uint   `json:"credit_count"`
}

type GetCurrencyResponse struct {
	Data   *CurrencyData `json:"data"`
	Status Status        `json:"status"`
}

type CurrencyData struct {
	ID     uint   `json:"id"`
	Symbol string `json:"symbol"`
	Slug   string `json:"slug"`
	Name   string `json:"name"`
}

func (e *CurrencyData) Currency() *currency.Currency {
	return &currency.Currency{
		ID:     e.ID,
		Symbol: e.Symbol,
		Slug:   e.Slug,
		Name:   e.Name,
	}
}

type DetailChartResponse struct {
	Data   *DetailChartData `json:"data"`
	Status Status           `json:"status"`
}

type DetailChartData struct {
	CurrencyID uint
	Points     *DetailChartPoints `json:"points"`
}

func (e *DetailChartData) PriceAndCapList() (*price_and_cap.PriceAndCapList, error) {
	if e == nil || e.Points == nil || len(*e.Points) == 0 || e.CurrencyID == 0 {
		return nil, apperror.ErrNotFound
	}
	var t string
	var point DetailChartPoint
	res := make(price_and_cap.PriceAndCapList, 0, len(*e.Points))
	for t, point = range *e.Points {
		i, err := strconv.ParseInt(t, 10, 64)
		if err != nil {
			return nil, err
		}

		if len(point.V) < 3 {
			continue
		}

		res = append(res, price_and_cap.PriceAndCap{
			CurrencyID:  e.CurrencyID,
			Price:       point.V[0],
			DailyVolume: point.V[1],
			Cap:         point.V[2],
			Ts:          time.Unix(i, 0),
		})
	}

	if len(res) == 0 {
		return nil, apperror.ErrNotFound
	}
	return &res, nil
}

type DetailChartPoints map[uint]DetailChartPoint

type DetailChartPoint struct {
	V []float64 `json:"v"`
	C []float64 `json:"c"`
}

type GetAnalyticsResponse struct {
	Data   *GetAnalyticsData `json:"data"`
	Status Status            `json:"status"`
}

type GetAnalyticsData struct {
	CurrencyID              uint
	AddressesByHoldings     *interface{}             `json:"addressesByHoldings"`
	HistoricalConcentration *HistoricalConcentration `json:"historicalConcentration"`
	AddressByTimeHeld       *interface{}             `json:"addressByTimeHeld"`
	AverageTransactionFees  *interface{}             `json:"averageTransactionFees"`
}

type HistoricalConcentration struct {
	HistoricalConcentrationDetails *[]HistoricalConcentrationDetailsPoint `json:"historicalConcentrationDetails"`
	HistoricalConcentrationAgg     *HistoricalConcentrationAgg            `json:"historicalConcentrationAgg"`
}

type HistoricalConcentrationAgg struct {
	WhalesPercent   float64 `json:"whalesPercent"`
	WhalesVolume    float64 `json:"whalesVolume"`
	OthersPercent   float64 `json:"othersPercent"`
	OthersVolume    float64 `json:"othersVolume"`
	InvestorPercent float64 `json:"investorPercent"`
	InvestorVolume  float64 `json:"investorVolume"`
	RetailPercent   float64 `json:"retailPercent"`
	RetailVolume    float64 `json:"retailVolume"`
}

type HistoricalConcentrationDetailsPoint struct {
	Date      string  `json:"date"`
	Whales    float64 `json:"whales"`
	Others    float64 `json:"others"`
	Retail    float64 `json:"retail"`
	Investors float64 `json:"investors"`
}

func (e *GetAnalyticsData) ConcentrationList() (*concentration.ConcentrationList, error) {
	if e == nil || e.HistoricalConcentration == nil || e.HistoricalConcentration.HistoricalConcentrationDetails == nil || len(*e.HistoricalConcentration.HistoricalConcentrationDetails) == 0 || e.CurrencyID == 0 {
		return nil, apperror.ErrNotFound
	}
	return e.HistoricalConcentration.ConcentrationList(e.CurrencyID)
}

func (e *HistoricalConcentration) ConcentrationList(currencyID uint) (*concentration.ConcentrationList, error) {
	if e == nil || e.HistoricalConcentrationDetails == nil || len(*e.HistoricalConcentrationDetails) == 0 || currencyID == 0 {
		return nil, apperror.ErrNotFound
	}

	res := make(concentration.ConcentrationList, 0, len(*e.HistoricalConcentrationDetails))
	var item HistoricalConcentrationDetailsPoint
	var t time.Time
	var err error
	for _, item = range *e.HistoricalConcentrationDetails {
		if t, err = time.Parse(time.DateOnly, item.Date); err != nil {
			return nil, err
		}

		res = append(res, concentration.Concentration{
			CurrencyID: currencyID,
			Whales:     item.Whales,
			Investors:  item.Investors,
			Retail:     item.Retail,
			D:          t,
		})
	}
	return &res, nil
}
