package invoicesrv

import (
	"context"
	"invoice-service/internal/domain"
)

type Service interface {
	CreateInvoice(ctx context.Context, input InvoiceInput) (*domain.Invoice, error)
	GetInvoice(ctx context.Context, invoiceID string) (*domain.Invoice, []int, error)
	GetIssuer(ctx context.Context, issuerID string) (*domain.Issuer, error)
	FindInvestors(ctx context.Context) ([]domain.Investor, error)
	PlaceBid(ctx context.Context, invoiceID string, investorID int, bidAmount float64) error
	FindTrades(ctx context.Context, tradeStatus *string) ([]domain.Trade, error)
	ApproveOrRejectTrade(ctx context.Context, tradeID string, isApproved bool) error
}

var _ Service = (*invoiceService)(nil)

type invoiceService struct {
	databaseRepository domain.DatabaseRepository
}

func NewService(databaseRepository domain.DatabaseRepository) *invoiceService {
	return &invoiceService{
		databaseRepository: databaseRepository,
	}
}
