package db

import (
	"context"
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"invoice-service/internal/domain"
	"invoice-service/pkg/utils"
	"testing"
)

func TestFindInvestors_ReturnsTrades(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)

	// Expected results
	var expectedInvestor = domain.NewInvestor(1, "some name", utils.PFloat(500.0))

	// Mock the expected rows and their data
	var rows = sqlmock.NewRows([]string{"id", "name", "available_funds"}).
		AddRow(1, "some name", 500.0)
	dep.sqlMock.ExpectQuery("SELECT (.+) FROM investors").
		WithArgs().
		WillReturnRows(rows).
		WillReturnError(nil)

	// Execute the function
	res, err := repo.FindInvestors(context.Background())

	// Assert the results
	assert.NoError(t, err)
	assert.Equal(t, []domain.Investor{*expectedInvestor}, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}

func TestFindInvestors_ReturnsError(t *testing.T) {
	// Initialize dependencies & Repository
	var dep = makeDependencies(t)
	var repo = initRepository(dep)


	// Mock the expected rows and their data
	dep.sqlMock.ExpectQuery("SELECT (.+) FROM investors").
		WithArgs().
		WillReturnRows().
		WillReturnError(sql.ErrNoRows)

	// Execute the function
	res, err := repo.FindInvestors(context.Background())

	// Assert the results
	assert.EqualError(t, err, "no investors found")
	assert.Nil(t, res)

	// Ensure all expectations were met
	assert.NoError(t, dep.sqlMock.ExpectationsWereMet())
}
