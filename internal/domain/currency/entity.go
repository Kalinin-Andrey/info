package currency

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
