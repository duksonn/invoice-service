package invoicesrv

import (
	"context"
	"invoice-service/internal/domain"
)

func (s *invoiceService) FindInvestors(ctx context.Context) ([]domain.Investor, error) {
	investors, err := s.databaseRepository.FindInvestors(ctx)
	if err != nil {
		return nil, err
	}

	return investors, nil
}
