package invoicehdlr

import (
	"errors"
	"invoice-service/pkg/server"
	"net/http"
)

type getInvoiceRequest struct {
	InvoiceID string
}

func (h *handler) GetInvoice(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	/** Build request DTO */
	var req = getInvoiceRequest{
		InvoiceID: server.GetStringFromPath(r, "invoice_id", ""),
	}

	/** Validate request */
	var err = req.validate()
	if err != nil {
		server.BadRequest(w, r, "BAD_REQUEST", err.Error())
		return
	}

	/** Service */
	invoice, investorsIDs, err := h.service.GetInvoice(ctx, req.InvoiceID)
	if err != nil {
		server.InternalServerError(w, r, err)
		return
	}

	server.OK(w, r, buildInvoiceResponse(invoice, investorsIDs))
}

func (r *getInvoiceRequest) validate() error {
	if r.InvoiceID == "" {
		return errors.New("invoice_id param is required")
	}
	return nil
}
