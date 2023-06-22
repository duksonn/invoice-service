package invoicehdlr

import (
	"invoice-service/pkg/server"
	"net/http"
)

func (h *handler) PostInvoice(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	/** Build request */
	req, err := newPostInvoiceRequest(r)
	if err != nil {
		server.BadRequest(w, r, "BAD_REQUEST", err.Error())
		return
	}

	/** Validate request */
	err = req.validate()
	if err != nil {
		server.BadRequest(w, r, "BAD_REQUEST", err.Error())
		return
	}

	/** Service */
	res, err := h.service.CreateInvoice(ctx, *req.Body.toInvoiceInput())
	if err != nil {
		server.InternalServerError(w, r, err)
		return
	}

	server.OK(w, r, buildInvoiceResponse(res, nil))
}
