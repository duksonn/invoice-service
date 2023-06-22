package http

import (
	"github.com/gorilla/mux"
	"invoice-service/cmd/dependencies"
	"invoice-service/cmd/http/handler/invoicehdlr"
	"net/http"
)

func routes(router mux.Router, dep *dependencies.Dependencies) *mux.Router {
	/** Handlers */
	var invoiceHdlr = invoicehdlr.NewInvoiceHandler(dep.InvoiceService)

	/** Routes */
	router.HandleFunc("/v1/invoice/{invoice_id}", invoiceHdlr.GetInvoice).Methods(http.MethodGet)
	router.HandleFunc("/v1/invoice", invoiceHdlr.PostInvoice).Methods(http.MethodPost)
	router.HandleFunc("/v1/issuer/{issuer_id}", invoiceHdlr.GetIssuer).Methods(http.MethodGet)
	router.HandleFunc("/v1/investor", invoiceHdlr.GetInvestors).Methods(http.MethodGet)
	router.HandleFunc("/v1/bid/place", invoiceHdlr.PostBid).Methods(http.MethodPost)
	router.HandleFunc("/v1/trade", invoiceHdlr.GetTrades).Methods(http.MethodGet)
	router.HandleFunc("/v1/trade/{trade_id}", invoiceHdlr.GetTrades).Methods(http.MethodGet)
	router.HandleFunc("/v1/trade/{trade_id}", invoiceHdlr.PutTrade).Methods(http.MethodPut)

	return &router
}
