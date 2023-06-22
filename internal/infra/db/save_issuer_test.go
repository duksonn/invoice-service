package db

import (
	"context"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"invoice-service/internal/domain"
	"invoice-service/pkg/utils"
	"testing"
)

func TestSaveIssuer_ValidIssuerExists_ReturnsIssuer(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedIssuer = domain.NewIssuer(1, "some name", utils.PFloat(500.0))

	// Mock the expected rows and their data
	dep.sqlMock.ExpectExec("INSERT INTO issuers(_*)").
		WithArgs().
		WillReturnResult(sqlmock.NewResult(1, 1)).
		WillReturnError(nil)

	// Execute the function
	res, err := repo.SaveIssuer(context.Background(), expectedIssuer)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedIssuer, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}

func TestSaveIssuer_ValidIssuerExists_ReturnsError(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedIssuer = domain.NewIssuer(1, "some name", utils.PFloat(500.0))

	// Mock the expected rows and their data
	dep.sqlMock.ExpectExec("INSERT INTO issuers(_*)").
		WithArgs().
		WillReturnError(sqlmock.ErrCancelled)

	// Execute the function
	res, err := repo.SaveIssuer(context.Background(), expectedIssuer)

	// Assert the results
	assert.EqualError(t, err, "canceling query due to user request")
	assert.Nil(t, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}
