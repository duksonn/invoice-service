package db

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"invoice-service/internal/domain"
	"testing"
)

func TestSaveBid_ValidBidExists_ReturnsBid(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedBid = domain.NewBid("id", 1, "invoiceID", 500.0, "some date")

	// Mock the expected rows and their data
	dep.sqlMock.ExpectExec("INSERT INTO bids(_*)").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)

	// Execute the function
	res, err := repo.SaveBid(context.Background(), expectedBid)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedBid, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}

func TestSaveBid_ValidBidExists_ReturnsError(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Mock the expected rows and their data
	dep.sqlMock.ExpectExec("INSERT INTO bids(_*)").
		WithArgs().
		WillReturnResult(nil).
		WillReturnError(sqlmock.ErrCancelled)

	// Execute the function
	var expectedBid = domain.NewBid("id", 1, "invoiceID", 500.0, "some date")
	res, err := repo.SaveBid(context.Background(), expectedBid)

	// Assert the results
	assert.EqualError(t, err, "canceling query due to user request")
	assert.Nil(t, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}
