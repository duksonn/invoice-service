package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"invoice-service/internal/domain"
	"invoice-service/internal/infra/db/dto"
)

func (m *DatabaseRepository) FindTrades(ctx context.Context, tradeStatus *string) ([]domain.Trade, error) {
	var query string
	var err error

	var retrievedTrades []dto.TradeDTO
	if tradeStatus == nil || *tradeStatus == "" {
		query = fmt.Sprintf("SELECT * FROM %s", TRADE_TABLE)
		err = m.client.Select(&retrievedTrades, query)
	} else {
		query = fmt.Sprintf("SELECT * FROM %s WHERE trade_status=$1", TRADE_TABLE)
		err = m.client.Select(&retrievedTrades, query, tradeStatus)
	}
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("no trades found"))
		}
		return nil, err
	}

	var trades []domain.Trade
	for _, t := range retrievedTrades {
		trades = append(trades, *t.ToTradeDomain())
	}

	return trades, nil
}
