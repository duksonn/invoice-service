package db

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	"invoice-service/internal/domain"
)

type DatabaseRepository struct {
	client *sqlx.DB
}

type DBConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	DbName   string `json:"db_name"`
	Schema   string `json:"schema"`
}

const (
	INVOICE_TABLE  = "invoices"
	ISSUER_TABLE   = "issuers"
	INVESTOR_TABLE = "investors"
	BID_TABLE      = "bids"
	TRADE_TABLE    = "trades"
)

var _ domain.DatabaseRepository = (*DatabaseRepository)(nil)

func NewDatabaseRepository(config DBConfig) *DatabaseRepository {
	var conn = fmt.Sprintf("postgres://%v:%v@%v:%v/%v?sslmode=disable&search_path=%v", config.User, config.Password, config.Host, config.Port, config.DbName, config.Schema)
	db, err := sqlx.Connect("postgres", conn)
	if err != nil {
		panic(err)
	}

	pingErr := db.Ping()
	if pingErr != nil {
		fmt.Println("cannot connect to database")
		panic(err)
	}
	fmt.Println("database connected!")

	return &DatabaseRepository{
		client: db,
	}
}
