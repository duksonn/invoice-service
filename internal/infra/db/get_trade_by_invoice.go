package db

import (
	"context"
	"fmt"
	"invoice-service/internal/domain"
	"invoice-service/internal/infra/db/dto"
)

func (m *DatabaseRepository) GetTradeByInvoice(ctx context.Context, invoiceID string) (*domain.Trade, error) {
	var query = fmt.Sprintf("SELECT * FROM %s WHERE invoice_id=$1", TRADE_TABLE)

	var retrievedTrade dto.TradeDTO
	var err = m.client.Get(&retrievedTrade, query, invoiceID)
	if err != nil {
		return nil, err
	}

	return retrievedTrade.ToTradeDomain(), nil
}
