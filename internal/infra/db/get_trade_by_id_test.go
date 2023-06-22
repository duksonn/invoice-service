package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"invoice-service/internal/domain"
	"testing"
)

func TestGetTradeByID_ValidTradeIDExists_ReturnsTrade(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedTradeID = "123"
	var expectedTrade = domain.NewTrade("123", "invoiceID", []int{1, 2}, "OPEN", "some date", nil)

	// Mock the expected rows and their data
	var rows = sqlmock.NewRows([]string{"id", "invoice_id", "investors_ids", "trade_status", "created_at", "updated_at"}).
		AddRow("123", "invoiceID", pq.Int64Array{1, 2}, "OPEN", "some date", nil)
	dep.sqlMock.ExpectQuery("SELECT (.+) FROM trades WHERE id=\\$1").
		WithArgs(expectedTradeID).
		WillReturnRows(rows).
		WillReturnError(nil)

	// Execute the function
	res, err := repo.GetTradeByID(context.Background(), expectedTradeID)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedTrade, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}

func TestGetTradeByID_ValidTradeIDDoesNotExist_ReturnsError(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedTradeID = "1234"

	// Mock the expected rows and their data
	dep.sqlMock.ExpectQuery("SELECT (.+) FROM trades WHERE id=\\$1").
		WithArgs(expectedTradeID).
		WillReturnError(sql.ErrNoRows)

	// Execute the function
	res, err := repo.GetTradeByID(context.Background(), expectedTradeID)

	// Assert the results
	var expectedError = fmt.Sprintf("trade %s not exist in db", expectedTradeID)
	assert.EqualError(t, err, expectedError)
	assert.Nil(t, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}
