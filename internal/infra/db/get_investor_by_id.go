package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"invoice-service/internal/domain"
	"invoice-service/internal/infra/db/dto"
)

func (m *DatabaseRepository) GetInvestorByID(ctx context.Context, investorID int) (*domain.Investor, error) {
	var query = fmt.Sprintf("SELECT * FROM %s WHERE id=$1", INVESTOR_TABLE)
	var row = m.client.QueryRowx(query, investorID)

	var retrievedInvestor dto.InvestorDTO
	if err := row.Scan(
		&retrievedInvestor.ID,
		&retrievedInvestor.Name,
		&retrievedInvestor.AvailableFunds,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("investor %v not exist in db", investorID))
		}
		return nil, err
	}

	return retrievedInvestor.ToInvestorDomain(), nil
}
