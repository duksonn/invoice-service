package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"invoice-service/internal/domain"
	"invoice-service/internal/infra/db/dto"
)

func (m *DatabaseRepository) GetTradeByID(ctx context.Context, tradeID string) (*domain.Trade, error) {
	var query = fmt.Sprintf("SELECT * FROM %s WHERE id=$1", TRADE_TABLE)

	var retrievedTrade dto.TradeDTO
	var err = m.client.Get(&retrievedTrade, query, tradeID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("trade %v not exist in db", tradeID))
		}
		return nil, err
	}

	return retrievedTrade.ToTradeDomain(), nil
}
