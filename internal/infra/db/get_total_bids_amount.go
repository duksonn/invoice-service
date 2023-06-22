package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

func (m *DatabaseRepository) GetTotalBidsAmount(ctx context.Context, invoiceID string) (*float64, error) {
	var query = fmt.Sprintf("SELECT SUM(bid_amount) AS total_bids_amount FROM %v WHERE invoice_id=$1", BID_TABLE)
	var row = m.client.QueryRowx(query, invoiceID)

	var retrievedAmount float64
	if err := row.Scan(&retrievedAmount); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("no bids for invoice %v", invoiceID))
		}
		return nil, err
	}

	return &retrievedAmount, nil
}
