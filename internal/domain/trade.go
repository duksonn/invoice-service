package domain

import (
	"invoice-service/pkg/utils"
	"time"
)

type Trade struct {
	id           string
	invoiceID    string
	investorsIDs []int
	tradeStatus  string
	createdAt    string
	updatedAt    *string
}

func NewTrade(id string, invoiceID string, investorsIDs []int, tradeStatus string, createdAt string, updatedAt *string) *Trade {
	return &Trade{
		id:           id,
		invoiceID:    invoiceID,
		investorsIDs: investorsIDs,
		tradeStatus:  tradeStatus,
		createdAt:    createdAt,
		updatedAt:    updatedAt,
	}
}

func CreateTrade(invoiceID string, investorsIDs []int) *Trade {
	var id = utils.GenerateUuid()
	var now = utils.Now().Format(time.RFC3339)

	return NewTrade(
		id,
		invoiceID,
		investorsIDs,
		WAITING_APPROVAL.String(),
		now,
		nil,
	)
}

func (t *Trade) UpgradeTradeStatus(isApproved bool) *Trade {
	var now = utils.Now().Format(time.RFC3339)

	var tradeStatus = REJECTED
	if isApproved {
		tradeStatus = ACCEPTED
	}

	return NewTrade(
		t.id,
		t.invoiceID,
		t.investorsIDs,
		tradeStatus.String(),
		t.createdAt,
		&now,
	)
}

func (t *Trade) ID() string {
	return t.id
}

func (t *Trade) InvoiceID() string {
	return t.invoiceID
}

func (t *Trade) InvestorsIDs() []int {
	return t.investorsIDs
}

func (t *Trade) TradeStatus() string {
	return t.tradeStatus
}

func (t *Trade) CreatedAt() string {
	return t.createdAt
}

func (t *Trade) UpdatedAt() *string {
	return t.updatedAt
}
