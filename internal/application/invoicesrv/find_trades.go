package invoicesrv

import (
	"context"
	"errors"
	"invoice-service/internal/domain"
)

func (s *invoiceService) FindTrades(ctx context.Context, tradeStatus *string) ([]domain.Trade, error) {
	if tradeStatus != nil && *tradeStatus != "" {
		if *tradeStatus != domain.WAITING_APPROVAL.String() && *tradeStatus != domain.ACCEPTED.String() && *tradeStatus != domain.REJECTED.String() {
			return nil, errors.New("invalid trade status")
		}
	}

	trades, err := s.databaseRepository.FindTrades(ctx, tradeStatus)
	if err != nil {
		return nil, err
	}

	return trades, nil
}
