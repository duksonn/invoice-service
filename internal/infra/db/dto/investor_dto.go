package dto

import "invoice-service/internal/domain"

type InvestorDTO struct {
	ID             int
	Name           string
	AvailableFunds *float64
}

func (i *InvestorDTO) ToInvestorDomain() *domain.Investor {
	return domain.NewInvestor(i.ID, i.Name, i.AvailableFunds)
}

func FromInvestorDomain(i *domain.Investor) InvestorDTO {
	return InvestorDTO{
		ID:             i.ID(),
		Name:           i.Name(),
		AvailableFunds: i.AvailableFunds(),
	}
}
