package currency

import "time"

const ()

type Currency struct {
	ID             uint
	Symbol         string
	Slug           string
	Name           string
	IsForObserving bool
}

func (e *Currency) Validate() error {
	return nil
}

type ImportMaxTime struct {
	CurrencyID    uint
	PriceAndCap   *time.Time
	Concentration *time.Time
}
