package domain

import "context"

//go:generate mockgen -source=./repository.go -package=mocks -destination=../../mocks/mockgen/database_repository.go

// Posibility to create interface for each table
type DatabaseRepository interface {
	// Invoice
	SaveInvoice(ctx context.Context, invoice *Invoice) (*Invoice, error)
	GetInvoiceByID(ctx context.Context, invoiceID string) (*Invoice, error)

	// Issuer
	GetIssuerByID(ctx context.Context, issuerID int) (*Issuer, error)
	SaveIssuer(ctx context.Context, issuer *Issuer) (*Issuer, error)

	// Investor
	FindInvestors(ctx context.Context) ([]Investor, error)
	GetInvestorByID(ctx context.Context, investorID int) (*Investor, error)
	SaveInvestor(ctx context.Context, investor *Investor) (*Investor, error)

	// Bid
	SaveBid(ctx context.Context, bid *Bid) (*Bid, error)
	GetTotalBidsAmount(ctx context.Context, invoiceID string) (*float64, error)
	GetBidsFromInvoiceAndInvestor(ctx context.Context, invoiceID string, investorID *int) ([]Bid, error)

	// Trade
	SaveTrade(ctx context.Context, trade *Trade) (*Trade, error)
	FindTrades(ctx context.Context, tradeStatus *string) ([]Trade, error)
	GetTradeByID(ctx context.Context, tradeID string) (*Trade, error)
	GetTradeByInvoice(ctx context.Context, invoiceID string) (*Trade, error)
}
