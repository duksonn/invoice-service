package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"invoice-service/pkg/utils"
	"testing"
)

func TestGetTotalBidsAmount_ValidInvoiceIDExists_ReturnsAmount(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedInvoiceID = "123"
	var expectedTotalAmount = utils.PFloat(500.0)

	// Mock the expected rows and their data
	rows := sqlmock.NewRows([]string{"total_bids_amount"}).AddRow(500.0)
	dep.sqlMock.ExpectQuery("SELECT SUM(.+) FROM bids WHERE invoice_id=\\$1").
		WithArgs(expectedInvoiceID).
		WillReturnRows(rows).
		WillReturnError(nil)

	// Execute the function
	totalAmount, err := repo.GetTotalBidsAmount(context.Background(), expectedInvoiceID)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedTotalAmount, totalAmount)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}

func TestGetTotalBidsAmount_ValidInvoiceIDDoesNotExist_ReturnsError(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedInvoiceID = "123"

	// Mock the expected rows and their data
	dep.sqlMock.ExpectQuery("SELECT SUM(.+) FROM bids WHERE invoice_id=\\$1").
		WithArgs(expectedInvoiceID).
		WillReturnError(sql.ErrNoRows)

	// Execute the function
	res, err := repo.GetTotalBidsAmount(context.Background(), expectedInvoiceID)

	// Assert the results
	var expectedError = fmt.Sprintf("no bids for invoice %v", expectedInvoiceID)
	assert.EqualError(t, err, expectedError)
	assert.Nil(t, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}
