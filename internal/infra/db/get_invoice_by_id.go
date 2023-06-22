package db

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"invoice-service/internal/domain"
	"invoice-service/internal/infra/db/dto"
)

func (m *DatabaseRepository) GetInvoiceByID(ctx context.Context, invoiceID string) (*domain.Invoice, error) {
	var query = fmt.Sprintf("SELECT * FROM %s WHERE id=$1", INVOICE_TABLE)
	var row = m.client.QueryRowx(query, invoiceID)

	var retrievedInvoice dto.InvoiceDTO
	var itemsData []byte
	if err := row.Scan(
		&retrievedInvoice.ID,
		&retrievedInvoice.DueDate,
		&retrievedInvoice.AskingPrice,
		&retrievedInvoice.Status,
		&itemsData,
		&retrievedInvoice.CreatedAt,
		&retrievedInvoice.IssuerID,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New(fmt.Sprintf("invoice %v not exist in db", invoiceID))
		}
		return nil, err
	}

	var err = json.Unmarshal(itemsData, &retrievedInvoice.Items)
	if err != nil {
		return nil, err
	}

	invoice, err := retrievedInvoice.ToInvoiceDomain()
	if err != nil {
		return nil, err
	}

	return invoice, nil
}
