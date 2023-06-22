package db

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
	"github.com/stretchr/testify/assert"
	"testing"
)

type dependencies struct {
	sqlDB   *sql.DB
	sqlMock sqlmock.Sqlmock
}

func makeDependencies(t *testing.T) *dependencies {
	// Create a new mock database
	dbMock, mock, err := sqlmock.New()
	assert.NoError(t, err)

	return &dependencies{
		sqlDB:   dbMock,
		sqlMock: mock,
	}
}

func initRepository(dep *dependencies) *DatabaseRepository {
	// Create a new DBConfig and pass it to the NewDatabaseRepository function
	var config = DBConfig{
		Host: "localhost",
		Port: 5432,
	}
	var repo = NewDatabaseRepository(config)

	// Replace the repository's client with the mock database connection
	repo.client = sqlx.NewDb(dep.sqlDB, "sqlmock")

	return repo
}
