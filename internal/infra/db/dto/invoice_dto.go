package dto

import (
	"invoice-service/internal/domain"
)

type InvoiceDTO struct {
	ID          string
	DueDate     string
	AskingPrice float64
	Status      string
	Items       []ItemDTO
	CreatedAt   string
	IssuerID    int
}

func (i *InvoiceDTO) ToInvoiceDomain() (*domain.Invoice, error) {
	var items []domain.Item
	for _, i := range i.Items {
		items = append(items, *domain.NewItem(i.ID, i.Description, i.Price, i.Quantity))
	}

	invoice, err := domain.NewInvoice(
		i.ID,
		i.DueDate,
		i.AskingPrice,
		i.Status,
		items,
		i.CreatedAt,
		i.IssuerID,
	)
	if err != nil {
		return nil, err
	}

	return invoice, err
}

func FromInvoiceDomain(invoice *domain.Invoice) InvoiceDTO {
	var items []ItemDTO
	for _, i := range invoice.Items() {
		items = append(items, FromItemDomain(i))
	}

	return InvoiceDTO{
		ID:          invoice.ID(),
		DueDate:     invoice.DueDate(),
		AskingPrice: invoice.AskingPrice(),
		Status:      invoice.Status(),
		Items:       items,
		CreatedAt:   invoice.CreatedAt(),
		IssuerID:    invoice.IssuerID(),
	}
}
