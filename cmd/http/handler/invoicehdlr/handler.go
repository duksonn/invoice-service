package invoicehdlr

import "invoice-service/internal/application/invoicesrv"

type handler struct {
	service invoicesrv.Service
}

func NewInvoiceHandler(service invoicesrv.Service) *handler {
	return &handler{
		service: service,
	}
}
