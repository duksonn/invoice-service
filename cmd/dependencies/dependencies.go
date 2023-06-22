package dependencies

import (
	"invoice-service/cmd/config"
	"invoice-service/internal/application/invoicesrv"
	"invoice-service/internal/infra/db"
)

type Dependencies struct {
	InvoiceService invoicesrv.Service
}

func Init(config *config.Config) *Dependencies {
	/** Repositories */
	var dbRepository = db.NewDatabaseRepository(config.MarketplaceMySql)

	/** Services */
	var invoiceSrv = invoicesrv.NewService(dbRepository)

	return &Dependencies{
		InvoiceService: invoiceSrv,
	}
}
