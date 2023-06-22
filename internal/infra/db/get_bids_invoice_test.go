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

func TestGetBidsFromInvoiceAndInvestor_ValidInvoiceIDExists_InvestorIDNil_ReturnsBids(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedInvoiceID = "invoiceID"
	var expectedBid = domain.NewBid("id", 1, "invoiceID", 500.0, "some date")

	// Mock the expected rows and their data
	var rows = sqlmock.NewRows([]string{"id", "investor_id", "invoice_id", "bid_amount", "created_at"}).
		AddRow("id", 1, "invoiceID", 500.0, "some date")
	dep.sqlMock.ExpectQuery("SELECT (.+) FROM bids WHERE invoice_id=\\$1").
		WithArgs(expectedInvoiceID).
		WillReturnRows(rows).
		WillReturnError(nil)

	// Execute the function
	res, err := repo.GetBidsFromInvoiceAndInvestor(context.Background(), expectedInvoiceID, nil)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, []domain.Bid{*expectedBid}, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}

func TestGetBidsFromInvoiceAndInvestor_ValidInvoiceIDExists_ValidInvestorIDExists_ReturnsBids(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedInvoiceID = "invoiceID"
	var expectedInvestorID = 1
	var expectedBid = domain.NewBid("id", 1, "invoiceID", 500.0, "some date")

	// Mock the expected rows and their data
	var rows = sqlmock.NewRows([]string{"id", "investor_id", "invoice_id", "bid_amount", "created_at"}).
		AddRow("id", 1, "invoiceID", 500.0, "some date")
	dep.sqlMock.ExpectQuery("SELECT (.+) FROM bids WHERE invoice_id=\\$1 AND investor_id=\\$2").
		WithArgs(expectedInvoiceID, expectedInvestorID).
		WillReturnRows(rows).
		WillReturnError(nil)

	// Execute the function
	res, err := repo.GetBidsFromInvoiceAndInvestor(context.Background(), expectedInvoiceID, &expectedInvestorID)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, []domain.Bid{*expectedBid}, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}

func TestGetBidsFromInvoiceAndInvestor_ValidInvoiceIDExists_InvestorIDNil_ReturnsError(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedInvoiceID = "invoiceID"

	// Mock the expected rows and their data
	dep.sqlMock.ExpectQuery("SELECT (.+) FROM bids WHERE invoice_id=\\$1").
		WithArgs(expectedInvoiceID).
		WillReturnError(sql.ErrNoRows)

	// Execute the function
	res, err := repo.GetBidsFromInvoiceAndInvestor(context.Background(), expectedInvoiceID, nil)

	// Assert the results
	var expectedError = fmt.Sprintf("no bids found for %v", expectedInvoiceID)
	assert.EqualError(t, err, expectedError)
	assert.Nil(t, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}
