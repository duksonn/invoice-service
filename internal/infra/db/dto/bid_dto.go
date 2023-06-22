package dto

import "invoice-service/internal/domain"

type BidDTO struct {
	ID         string
	InvestorID int
	InvoiceID  string
	BidAmount  float64
	CreatedAt  string
}

func (b *BidDTO) ToBidDomain() *domain.Bid {
	return domain.NewBid(b.ID, b.InvestorID, b.InvoiceID, b.BidAmount, b.CreatedAt)
}

func FromBidDomain(b *domain.Bid) *BidDTO {
	return &BidDTO{
		ID:         b.ID(),
		InvestorID: b.InvestorID(),
		InvoiceID:  b.InvoiceID(),
		BidAmount:  b.BidAmount(),
		CreatedAt:  b.CreatedAt(),
	}
}
