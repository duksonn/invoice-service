package db

import (
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"golang.org/x/net/context"
	"invoice-service/internal/domain"
	"invoice-service/pkg/utils"
	"testing"
)

func TestGetIssuerByID_ValidIssuerIDExists_ReturnsIssuer(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedIssuerID = 1
	var expectedIssuer = domain.NewIssuer(1, "some name", utils.PFloat(500.0))

	// Mock the expected rows and their data
	var rows = sqlmock.NewRows([]string{"id", "company_name", "available_funds"}).
		AddRow(1, "some name", 500.0)
	dep.sqlMock.ExpectQuery("SELECT (.+) FROM issuers WHERE id=\\$1").
		WithArgs(expectedIssuerID).
		WillReturnRows(rows).
		WillReturnError(nil)

	// Execute the function
	res, err := repo.GetIssuerByID(context.Background(), expectedIssuerID)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedIssuer, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}

func TestGetIssuerByID_ValidIssuerIDDoesNotExist_ReturnsError(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedIssuerID = 1

	// Mock the expected rows and their data
	dep.sqlMock.ExpectQuery("SELECT (.+) FROM issuers WHERE id=\\$1").
		WithArgs(expectedIssuerID).
		WillReturnError(sql.ErrNoRows)

	// Execute the function
	res, err := repo.GetIssuerByID(context.Background(), expectedIssuerID)

	// Assert the results
	var expectedError = fmt.Sprintf("issuer %v not exist in db", expectedIssuerID)
	assert.EqualError(t, err, expectedError)
	assert.Nil(t, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}
