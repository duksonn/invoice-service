package dto

import "invoice-service/internal/domain"

type IssuerDTO struct {
	ID             int
	CompanyName    string
	AvailableFunds *float64
}

func (i *IssuerDTO) ToIssuerDomain() *domain.Issuer {
	return domain.NewIssuer(i.ID, i.CompanyName, i.AvailableFunds)
}

func FromIssuerDomain(i *domain.Issuer) *IssuerDTO {
	return &IssuerDTO{
		ID:             i.ID(),
		CompanyName:    i.CompanyName(),
		AvailableFunds: i.AvailableFunds(),
	}
}
