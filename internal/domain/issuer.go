package domain

type Issuer struct {
	id             int
	companyName    string
	availableFunds *float64
}

func NewIssuer(id int, companyName string, availableFunds *float64) *Issuer {
	return &Issuer{
		id:             id,
		companyName:    companyName,
		availableFunds: availableFunds,
	}
}

func (i *Issuer) UpdateBalance(availableFunds float64) {
	i.availableFunds = &availableFunds
}

func (i *Issuer) ID() int {
	return i.id
}

func (i *Issuer) CompanyName() string {
	return i.companyName
}

func (i *Issuer) AvailableFunds() *float64 {
	return i.availableFunds
}
