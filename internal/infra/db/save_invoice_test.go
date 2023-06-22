package db

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"invoice-service/internal/domain"
	"testing"
)

func TestSaveInvoice_ValidInvoiceExists_ReturnsInvoice(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	expectedInvoice, err := domain.NewInvoice(
		"id",
		"some date",
		500.0,
		"OPEN",
		[]domain.Item{*domain.NewItem("item_1", "some desc", 500.0, 1)},
		"some date",
		1,
	)

	// Mock the expected rows and their data
	dep.sqlMock.ExpectExec("INSERT INTO invoices(_*)").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)

	// Execute the function
	res, err := repo.SaveInvoice(context.Background(), expectedInvoice)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedInvoice, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}

func TestSaveInvoice_ValidInvoiceExists_ReturnsError(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	expectedInvoice, err := domain.NewInvoice(
		"id",
		"some date",
		500.0,
		"OPEN",
		[]domain.Item{*domain.NewItem("item_1", "some desc", 500.0, 1)},
		"some date",
		1,
	)

	// Mock the expected rows and their data
	dep.sqlMock.ExpectExec("INSERT INTO invoices(_*)").
		WithArgs(
			expectedInvoice.ID(),
			expectedInvoice.DueDate(),
			expectedInvoice.AskingPrice(),
			expectedInvoice.Status(),
			sqlmock.AnyArg(),
			expectedInvoice.CreatedAt(),
			expectedInvoice.IssuerID(),
		).
		WillReturnError(sqlmock.ErrCancelled)

	// Execute the function
	res, err := repo.SaveInvoice(context.Background(), expectedInvoice)

	// Assert the results
	assert.EqualError(t, err, "canceling query due to user request")
	assert.Nil(t, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}
