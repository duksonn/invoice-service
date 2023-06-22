package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"invoice-service/internal/domain"
	"invoice-service/internal/infra/db/dto"
)

func (m *DatabaseRepository) GetIssuerByID(ctx context.Context, issuerID int) (*domain.Issuer, error) {
	var query = fmt.Sprintf("SELECT * FROM %s WHERE id=$1", ISSUER_TABLE)
	var row = m.client.QueryRowx(query, issuerID)

	var retrievedIssuer dto.IssuerDTO
	if err := row.Scan(
		&retrievedIssuer.ID,
		&retrievedIssuer.CompanyName,
		&retrievedIssuer.AvailableFunds,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("issuer %v not exist in db", issuerID))
		}
		return nil, err
	}

	return retrievedIssuer.ToIssuerDomain(), nil
}
