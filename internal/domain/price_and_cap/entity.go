package price_and_cap

import "time"

const ()

type PriceAndCap struct {
	CurrencyID  uint
	Price       float64
	DailyVolume float64
	Cap         float64
	Ts          time.Time
}

func (e *PriceAndCap) Validate() error {
	return nil
}

type PriceAndCapList []PriceAndCap
