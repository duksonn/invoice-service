package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"invoice-service/internal/domain"
	"invoice-service/internal/infra/db/dto"
)

func (m *DatabaseRepository) FindInvestors(ctx context.Context) ([]domain.Investor, error) {
	var query = fmt.Sprintf("SELECT * FROM %s", INVESTOR_TABLE)
	rows, err := m.client.Queryx(query)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("no investors found"))
		}
		return nil, err
	}
	defer rows.Close()

	var investors []domain.Investor
	for rows.Next() {
		var retrievedInvestor dto.InvestorDTO
		if err := rows.Scan(
			&retrievedInvestor.ID,
			&retrievedInvestor.Name,
			&retrievedInvestor.AvailableFunds,
		); err != nil {
			return nil, err
		}
		investors = append(investors, *retrievedInvestor.ToInvestorDomain())
	}

	return investors, nil
}
