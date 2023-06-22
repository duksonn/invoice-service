package db

import (
	"context"
	"fmt"
	"invoice-service/internal/domain"
	"invoice-service/internal/infra/db/dto"
)

func (m *DatabaseRepository) SaveBid(ctx context.Context, bid *domain.Bid) (*domain.Bid, error) {
	var bidDTO = dto.FromBidDomain(bid)

	var query = fmt.Sprintf("INSERT INTO %s(id, investor_id, invoice_id, bid_amount, created_at) VALUES ('%v', %v, '%v', %v, '%v') "+
		"ON CONFLICT (id) DO UPDATE SET investor_id=EXCLUDED.investor_id, invoice_id=EXCLUDED.invoice_id, "+
		"bid_amount=EXCLUDED.bid_amount, created_at=EXCLUDED.created_at",
		BID_TABLE,
		bidDTO.ID,
		bidDTO.InvestorID,
		bidDTO.InvoiceID,
		bidDTO.BidAmount,
		bidDTO.CreatedAt,
	)
	_, err := m.client.Exec(query)
	if err != nil {
		return nil, err
	}

	return bid, nil
}
