package invoicehdlr

import (
	"github.com/golang/mock/gomock"
	"invoice-service/internal/application/invoicesrv"
	mocks "invoice-service/mocks/mockgen"
	"testing"
)

type dependencies struct {
	databaseRepository *mocks.MockDatabaseRepository
}

type errResponse struct {
	Messages []string `json:"messages"`
	Code     string   `json:"code"`
}

func makeDependencies(t *testing.T) *dependencies {
	return &dependencies{
		databaseRepository: mocks.NewMockDatabaseRepository(gomock.NewController(t)),
	}
}

func buildHandler(dep *dependencies) *handler {
	var invoiceSrv = invoicesrv.NewService(dep.databaseRepository)
	return NewInvoiceHandler(invoiceSrv)
}
