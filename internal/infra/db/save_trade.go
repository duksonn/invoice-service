package db

import (
	"context"
	"fmt"
	"invoice-service/internal/domain"
	"invoice-service/internal/infra/db/dto"
)

func (m *DatabaseRepository) SaveTrade(ctx context.Context, trade *domain.Trade) (*domain.Trade, error) {
	var tradeDTO = dto.FromTradeDomain(trade)

	var query = fmt.Sprintf("INSERT INTO %s(id, invoice_id, investors_ids, trade_status, created_at, updated_at) "+
		"VALUES ($1, $2, $3, $4, $5, $6) "+
		"ON CONFLICT (id) DO UPDATE SET invoice_id=$2, investors_ids=$3, trade_status=$4, created_at=$5, updated_at=$6",
		TRADE_TABLE,
	)
	_, err := m.client.Exec(
		query,
		tradeDTO.ID,
		tradeDTO.InvoiceID,
		tradeDTO.InvestorsIDs,
		tradeDTO.TradeStatus,
		tradeDTO.CreatedAt,
		tradeDTO.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	return trade, nil
}
