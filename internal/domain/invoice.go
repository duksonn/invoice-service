package domain

import (
	"invoice-service/pkg/utils"
	"time"
)

type Invoice struct {
	id          string
	dueDate     string
	askingPrice float64
	status      string
	items       []Item
	createdAt   string
	issuerID    int
}

func NewInvoice(
	id string,
	dueDate string,
	askingPrice float64,
	status string,
	items []Item,
	createdAt string,
	issuerID int,
) (*Invoice, error) {
	var invoice = &Invoice{
		id:          id,
		dueDate:     dueDate,
		askingPrice: askingPrice,
		status:      status,
		items:       items,
		createdAt:   createdAt,
		issuerID:    issuerID,
	}

	// Validations: Could make some validations here

	// Normalize: Could normalize data here before persisting in db

	return invoice, nil
}

func CreateInvoice(items []Item, issuerID int) (*Invoice, error) {
	var id = utils.GenerateUuid()
	var now = utils.Now().Format(time.RFC3339)
	var dueDate = utils.Now().AddDate(0, 0, 21).Format(time.RFC3339) // 3 weeks to due date
	var askingPrice = calculateTotalAmount(items)

	invoice, err := NewInvoice(
		id,
		dueDate,
		askingPrice,
		ACTIVE.String(),
		items,
		now,
		issuerID,
	)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}

func (i *Invoice) UpdateInvoice(status string) (*Invoice, error) {
	newInvoice, err := NewInvoice(
		i.id,
		i.dueDate,
		i.askingPrice,
		status,
		i.items,
		i.createdAt,
		i.issuerID,
	)
	if err != nil {
		return nil, err
	}

	return newInvoice, nil
}

func calculateTotalAmount(items []Item) float64 {
	var totalAmount float64
	for _, i := range items {
		totalAmount = totalAmount + (i.Price() * float64(i.Quantity()))
	}

	return totalAmount
}

func (i *Invoice) ID() string {
	return i.id
}

func (i *Invoice) DueDate() string {
	return i.dueDate
}

func (i *Invoice) AskingPrice() float64 {
	return i.askingPrice
}

func (i *Invoice) Status() string {
	return i.status
}

func (i *Invoice) Items() []Item {
	return i.items
}

func (i *Invoice) CreatedAt() string {
	return i.createdAt
}

func (i *Invoice) IssuerID() int {
	return i.issuerID
}
