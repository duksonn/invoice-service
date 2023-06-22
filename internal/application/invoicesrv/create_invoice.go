package invoicesrv

import (
	"context"
	"invoice-service/internal/domain"
)

type InvoiceInput struct {
	Items    []ItemInput
	IssuerID int
}

type ItemInput struct {
	ID          string
	Description string
	Price       float64
	Quantity    int
}

func (s *invoiceService) CreateInvoice(ctx context.Context, input InvoiceInput) (*domain.Invoice, error) {
	/** Convert items to domain */
	var items []domain.Item
	for _, i := range input.Items {
		items = append(items, *domain.NewItem(i.ID, i.Description, i.Price, i.Quantity))
	}

	/** Create invoice domain */
	invoice, err := domain.CreateInvoice(items, input.IssuerID)
	if err != nil {
		return nil, err
	}

	/** Save invoice in DB */
	res, err := s.databaseRepository.SaveInvoice(ctx, invoice)
	if err != nil {
		return nil, err
	}

	return res, nil
}
