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

func (l *PriceAndCapList) Slice() *[]PriceAndCap {
	if l == nil {
		return nil
	}
	res := []PriceAndCap(*l)
	return &res
}

func (l *PriceAndCapList) MaxTime() *time.Time {
	if l == nil || len(*l) == 0 {
		return nil
	}
	max := (*l)[0].Ts
	var item PriceAndCap
	for _, item = range *l {
		if item.Ts.After(max) {
			max = item.Ts
		}
	}
	return &max
}

type PriceAndCapMap map[uint]PriceAndCapList
