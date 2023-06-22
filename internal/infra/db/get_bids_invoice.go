package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx"
	"invoice-service/internal/domain"
	"invoice-service/internal/infra/db/dto"
)

func (m *DatabaseRepository) GetBidsFromInvoiceAndInvestor(ctx context.Context, invoiceID string, investorID *int) ([]domain.Bid, error) {
	var rows *sqlx.Rows
	var err error

	if investorID == nil || *investorID == 0 {
		var query = fmt.Sprintf("SELECT * FROM %v WHERE invoice_id=$1", BID_TABLE)
		rows, err = m.client.Queryx(query, invoiceID)
	} else {
		var query = fmt.Sprintf("SELECT * FROM %v WHERE invoice_id=$1 AND investor_id=$2", BID_TABLE)
		rows, err = m.client.Queryx(query, invoiceID, investorID)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("no bids found for %v", invoiceID))
		}
		return nil, err
	}

	var bids []domain.Bid
	for rows.Next() {
		var retrievedBid dto.BidDTO
		if err := rows.Scan(
			&retrievedBid.ID,
			&retrievedBid.InvestorID,
			&retrievedBid.InvoiceID,
			&retrievedBid.BidAmount,
			&retrievedBid.CreatedAt,
		); err != nil {
			return nil, err
		}
		bids = append(bids, *retrievedBid.ToBidDomain())
	}

	return bids, nil
}
