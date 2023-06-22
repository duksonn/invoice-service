package config

import "invoice-service/internal/infra/db"

type Config struct {
	Server           ServerConfig `json:"server"`
	MarketplaceMySql db.DBConfig  `json:"postgre_db"`
}

type ServerConfig struct {
	Port int `json:"port"`
}
