package db

import (
	"context"
	"fmt"
	"invoice-service/internal/domain"
	"invoice-service/internal/infra/db/dto"
)

func (m *DatabaseRepository) SaveInvestor(ctx context.Context, investor *domain.Investor) (*domain.Investor, error) {
	var investorDTO = dto.FromInvestorDomain(investor)

	var query = fmt.Sprintf("INSERT INTO %s(id, name, available_funds) VALUES (%v, '%v', %v) "+
		"ON CONFLICT (id) DO UPDATE SET name=EXCLUDED.name, available_funds=EXCLUDED.available_funds",
		INVESTOR_TABLE,
		investorDTO.ID,
		investorDTO.Name,
		*investorDTO.AvailableFunds,
	)
	_, err := m.client.Exec(query)
	if err != nil {
		return nil, err
	}

	return investor, nil
}
