package currency

import "time"

const ()

type ImportMaxTime struct {
	CurrencyID    uint
	PriceAndCap   *time.Time
	Concentration *time.Time
}

type Currency struct {
	ID                uint
	Symbol            string
	Slug              string
	Name              string
	IsForObserving    bool
	CirculatingSupply uint
	TotalSupply       uint
	MaxSupply         uint
	LatestPrice       float64
	CmcRank           uint
	AddedAt           uint
	Platform          *CurrencyPlatform
}

type CurrencyPlatform struct {
	ID           uint
	Symbol       string
	Slug         string
	Name         string
	TokenAddress string
}

func (e *Currency) Validate() error {
	return nil
}

type CurrencyList []Currency

func (l *CurrencyList) IDs() *[]uint {
	if l == nil {
		return nil
	}
	res := make([]uint, 0, len(*l))
	var item Currency
	for _, item = range *l {
		res = append(res, item.ID)
	}
	return &res
}

type CurrencyMap map[uint]Currency
