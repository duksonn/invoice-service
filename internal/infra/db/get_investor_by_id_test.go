package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"invoice-service/internal/domain"
	"invoice-service/pkg/utils"
	"testing"
)

func TestGetInvestorByID_ValidInvestorIDExists_ReturnsInvestor(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedInvestorID = 1
	var expectedInvestor = domain.NewInvestor(1, "some name", utils.PFloat(500.0))

	// Mock the expected rows and their data
	var rows = sqlmock.NewRows([]string{"id", "name", "available_funds"}).
		AddRow(1, "some name", 500.0)
	dep.sqlMock.ExpectQuery("SELECT (.+) FROM investors WHERE id=\\$1").
		WithArgs(expectedInvestorID).
		WillReturnRows(rows).
		WillReturnError(nil)

	// Execute the function
	res, err := repo.GetInvestorByID(context.Background(), expectedInvestorID)

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, expectedInvestor, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}

func TestGetInvestorByID_ValidInvestorIDDoesNotExist_ReturnsError(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedInvestorID = 1

	// Mock the expected rows and their data
	dep.sqlMock.ExpectQuery("SELECT (.+) FROM investors WHERE id=\\$1").
		WithArgs(expectedInvestorID).
		WillReturnError(sql.ErrNoRows)

	// Execute the function
	res, err := repo.GetInvestorByID(context.Background(), expectedInvestorID)

	// Assert the results
	var expectedError = fmt.Sprintf("investor %v not exist in db", expectedInvestorID)
	assert.EqualError(t, err, expectedError)
	assert.Nil(t, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}
