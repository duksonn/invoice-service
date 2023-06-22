package invoicesrv

import (
	"context"
	"database/sql"
	"invoice-service/internal/domain"
)

func (s *invoiceService) GetInvoice(ctx context.Context, invoiceID string) (*domain.Invoice, []int, error) {
	// Retrieve the invoice
	invoice, err := s.databaseRepository.GetInvoiceByID(ctx, invoiceID)
	if err != nil {
		return nil, nil, err
	}

	// Check if the invoice has been purchased
	purchased, trade, err := s.checkInvoicePurchase(ctx, invoiceID)
	if err != nil {
		return nil, nil, err
	}

	if purchased {
		return invoice, trade.InvestorsIDs(), nil
	}

	return invoice, nil, nil
}

func (s *invoiceService) checkInvoicePurchase(ctx context.Context, invoiceID string) (bool, *domain.Trade, error) {
	trade, err := s.databaseRepository.GetTradeByInvoice(ctx, invoiceID)
	if err != nil {
		if err == sql.ErrNoRows {
			return false, nil, nil
		}
		return false, nil, err
	}

	return true, trade, nil
}
