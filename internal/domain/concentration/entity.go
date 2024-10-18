package concentration

import "time"

const ()

type Concentration struct {
	CurrencyID uint
	Whales     float64
	Investors  float64
	Retail     float64
	D          time.Time
}

func (e *Concentration) Validate() error {
	return nil
}

type ConcentrationList []Concentration

func (l *ConcentrationList) Slice() *[]Concentration {
	if l == nil {
		return nil
	}
	res := []Concentration(*l)
	return &res
}

type ConcentrationMap map[uint]ConcentrationList
