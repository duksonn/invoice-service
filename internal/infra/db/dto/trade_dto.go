package dto

import (
	"github.com/lib/pq"
	"invoice-service/internal/domain"
)

// In this case we need the tag 'db' cause we are using Get and Select from sqlx
type TradeDTO struct {
	ID           string        `db:"id"`
	InvoiceID    string        `db:"invoice_id"`
	InvestorsIDs pq.Int64Array `db:"investors_ids"`
	TradeStatus  string        `db:"trade_status"`
	CreatedAt    string        `db:"created_at"`
	UpdatedAt    *string       `db:"updated_at"`
}

func (t *TradeDTO) ToTradeDomain() *domain.Trade {
	var investorsIDs []int
	for _, i := range t.InvestorsIDs {
		investorsIDs = append(investorsIDs, int(i))
	}

	return domain.NewTrade(t.ID, t.InvoiceID, investorsIDs, t.TradeStatus, t.CreatedAt, t.UpdatedAt)
}

func FromTradeDomain(t *domain.Trade) *TradeDTO {
	var investorsIDs []int64
	for _, i := range t.InvestorsIDs() {
		investorsIDs = append(investorsIDs, int64(i))
	}

	return &TradeDTO{
		ID:           t.ID(),
		InvoiceID:    t.InvoiceID(),
		InvestorsIDs: investorsIDs,
		TradeStatus:  t.TradeStatus(),
		CreatedAt:    t.CreatedAt(),
		UpdatedAt:    t.UpdatedAt(),
	}
}
