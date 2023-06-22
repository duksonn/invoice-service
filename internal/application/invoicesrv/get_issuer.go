package invoicesrv

import (
	"context"
	"invoice-service/internal/domain"
	"strconv"
)

func (s *invoiceService) GetIssuer(ctx context.Context, issuerID string) (*domain.Issuer, error) {
	id, err := strconv.Atoi(issuerID)
	if err != nil {
		return nil, err
	}

	issuer, err := s.databaseRepository.GetIssuerByID(ctx, id)
	if err != nil {
		return nil, err
	}

	return issuer, nil
}
