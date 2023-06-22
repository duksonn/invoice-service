package domain

import (
	"invoice-service/pkg/utils"
	"time"
)

type Bid struct {
	id         string
	investorID int
	invoiceID  string
	bidAmount  float64
	createdAt  string
}

func NewBid(id string, investorID int, invoiceID string, bidAmount float64, createdAt string) *Bid {
	return &Bid{
		id:         id,
		investorID: investorID,
		invoiceID:  invoiceID,
		bidAmount:  bidAmount,
		createdAt:  createdAt,
	}
}

func CreateBid(investorID int, invoiceID string, bidAmount float64) *Bid {
	var id = utils.GenerateUuid()
	var now = utils.Now().Format(time.RFC3339)

	return NewBid(id, investorID, invoiceID, bidAmount, now)
}

func (b *Bid) ID() string {
	return b.id
}

func (b *Bid) InvestorID() int {
	return b.investorID
}

func (b *Bid) InvoiceID() string {
	return b.invoiceID
}

func (b *Bid) BidAmount() float64 {
	return b.bidAmount
}

func (b *Bid) CreatedAt() string {
	return b.createdAt
}
