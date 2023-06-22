package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"invoice-service/internal/domain"
	"testing"
)

func TestGetInvoiceByID_ValidInvoiceIDExists_ReturnsInvoice(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedInvoiceID = "id"
	expectedInvoice, _ := domain.NewInvoice(
		"id",
		"some date",
		500.0,
		"OPEN",
		[]domain.Item{*domain.NewItem("item_1", "some desc", 500.0, 1)},
		"some date",
		1,
	)

	// Mock the expected rows and their data
	var rows = sqlmock.NewRows([]string{"id", "due_date", "asking_price", "status", "items", "created_at", "issuer_id"}).
		AddRow("id", "some date", 500.0, "OPEN",
			"[{\"id\": \"item_1\", \"price\": 500.0, \"quantity\": 1, \"description\": \"some desc\"}]",
			"some date", 1)
	dep.sqlMock.ExpectQuery("SELECT (.+) FROM invoices WHERE id=\\$1").
		WithArgs(expectedInvoiceID).
		WillReturnRows(rows).
		WillReturnError(nil)

	// Execute the function
	res, err := repo.GetInvoiceByID(context.Background(), expectedInvoiceID)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedInvoice, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}

func TestGetInvoiceByID_ValidInvoiceIDDoesNotExist_ReturnsError(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedInvoiceID = "id"

	// Mock the expected rows and their data
	dep.sqlMock.ExpectQuery("SELECT (.+) FROM invoices WHERE id=\\$1").
		WithArgs(expectedInvoiceID).
		WillReturnError(sql.ErrNoRows)

	// Execute the function
	res, err := repo.GetInvoiceByID(context.Background(), expectedInvoiceID)

	// Assert the results
	var expectedError = fmt.Sprintf("invoice %v not exist in db", expectedInvoiceID)
	assert.EqualError(t, err, expectedError)
	assert.Nil(t, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}
