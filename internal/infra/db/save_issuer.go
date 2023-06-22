package db

import (
	"context"
	"fmt"
	"invoice-service/internal/domain"
	"invoice-service/internal/infra/db/dto"
)

func (m *DatabaseRepository) SaveIssuer(ctx context.Context, issuer *domain.Issuer) (*domain.Issuer, error) {
	var issuerDTO = dto.FromIssuerDomain(issuer)

	var query = fmt.Sprintf("INSERT INTO %s(id, company_name, available_funds) VALUES (%v, '%v', %v) "+
		"ON CONFLICT (id) DO UPDATE SET company_name=EXCLUDED.company_name, available_funds=EXCLUDED.available_funds",
		ISSUER_TABLE,
		issuerDTO.ID,
		issuerDTO.CompanyName,
		*issuerDTO.AvailableFunds,
	)
	_, err := m.client.Exec(query)
	if err != nil {
		return nil, err
	}

	return issuer, nil
}
