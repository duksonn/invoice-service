package db

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"invoice-service/internal/domain"
	"testing"
)

func TestSaveTrade_ValidTradeExists_ReturnsTrade(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedTrade = domain.NewTrade("123", "invoiceID", []int{1, 2}, "OPEN", "some date", nil)

	// Mock the expected rows and their data
	dep.sqlMock.ExpectExec("INSERT INTO trades(_*)").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)

	// Execute the function
	res, err := repo.SaveTrade(context.Background(), expectedTrade)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedTrade, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}

func TestSaveTrade_ValidTradeExists_ReturnsError(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedTrade = domain.NewTrade("123", "invoiceID", []int{1, 2}, "OPEN", "some date", nil)

	// Mock the expected rows and their data
	dep.sqlMock.ExpectExec("INSERT INTO trades(_*)").
		WithArgs().
		WillReturnError(sqlmock.ErrCancelled)

	// Execute the function
	res, err := repo.SaveTrade(context.Background(), expectedTrade)

	// Assert the results
	assert.EqualError(t, err, "canceling query due to user request")
	assert.Nil(t, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}
