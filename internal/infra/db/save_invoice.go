package db

import (
	"context"
	"encoding/json"
	"fmt"
	"invoice-service/internal/domain"
	"invoice-service/internal/infra/db/dto"
)

func (m *DatabaseRepository) SaveInvoice(ctx context.Context, invoice *domain.Invoice) (*domain.Invoice, error) {
	var invoiceDTO = dto.FromInvoiceDomain(invoice)

	itemsJson, err := json.Marshal(invoiceDTO.Items)
	if err != nil {
		return nil, err
	}

	var query = fmt.Sprintf("INSERT INTO %s(id, due_date, asking_price, status, items, created_at, issuer_id) "+
		"VALUES ($1, $2, $3, $4, $5, $6, $7) ON CONFLICT (id) DO UPDATE SET due_date=$2, "+
		"asking_price=$3, status=$4, items=$5, created_at=$6, issuer_id=$7", INVOICE_TABLE)
	_, err = m.client.Exec(query, invoiceDTO.ID, invoiceDTO.DueDate, invoiceDTO.AskingPrice, invoiceDTO.Status, itemsJson, invoiceDTO.CreatedAt, invoiceDTO.IssuerID)
	if err != nil {
		return nil, err
	}

	return invoice, nil
}
