package domain

type Investor struct {
	id             int
	name           string
	availableFunds *float64
}

func NewInvestor(id int, name string, availableFunds *float64) *Investor {
	return &Investor{
		id:             id,
		name:           name,
		availableFunds: availableFunds,
	}
}

func (i *Investor) UpdateBalance(availableFunds float64) {
	i.availableFunds = &availableFunds
}

func (i *Investor) ID() int {
	return i.id
}

func (i *Investor) Name() string {
	return i.name
}

func (i *Investor) AvailableFunds() *float64 {
	return i.availableFunds
}
