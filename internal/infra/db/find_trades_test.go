package db

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/lib/pq"
	"github.com/stretchr/testify/assert"
	"invoice-service/internal/domain"
	"testing"
)

func TestFindTrades_NoStatus_ReturnsTrades(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedTrade = domain.NewTrade("123", "invoiceID", []int{1, 2}, "OPEN", "some date", nil)

	var rows = sqlmock.NewRows([]string{"id", "invoice_id", "investors_ids", "trade_status", "created_at", "updated_at"}).
		AddRow("123", "invoiceID", pq.Int64Array{1, 2}, "OPEN", "some date", nil)
	dep.sqlMock.ExpectQuery("SELECT (.+) FROM trades").
		WithArgs().
		WillReturnRows(rows).
		WillReturnError(nil)

	// Execute the function
	res, err := repo.FindTrades(context.Background(), nil)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, []domain.Trade{*expectedTrade}, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}

func TestFindTrades_ValidStatus_ReturnsTrades(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedStatus = "WAITING_APPROVAL"
	var expectedTrade = domain.NewTrade("123", "invoiceID", []int{1, 2}, "WAITING_APPROVAL", "some date", nil)

	var rows = sqlmock.NewRows([]string{"id", "invoice_id", "investors_ids", "trade_status", "created_at", "updated_at"}).
		AddRow("123", "invoiceID", pq.Int64Array{1, 2}, "WAITING_APPROVAL", "some date", nil)
	dep.sqlMock.ExpectQuery("SELECT (.+) FROM trades WHERE trade_status=\\$1").
		WithArgs(expectedStatus).
		WillReturnRows(rows).
		WillReturnError(nil)

	// Execute the function
	res, err := repo.FindTrades(context.Background(), &expectedStatus)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, []domain.Trade{*expectedTrade}, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}

func TestFindTrades_NoStatus_ReturnsError(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	dep.sqlMock.ExpectQuery("SELECT (.+) FROM trades").
		WithArgs().
		WillReturnError(sql.ErrNoRows)

	// Execute the function
	res, err := repo.FindTrades(context.Background(), nil)

	// Assert the results
	assert.EqualError(t, err, "no trades found")
	assert.Nil(t, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}

