package invoicehdlr

import (
	"invoice-service/pkg/server"
	"net/http"
)

func (h *handler) PostBid(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	/** Build request */
	req, err := newPostBidRequest(r)
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
	err = h.service.PlaceBid(ctx, req.Body.InvoiceID, req.Body.InvestorID, req.Body.BidAmount)
	if err != nil {
		server.InternalServerError(w, r, err)
		return
	}

	server.OKNoContent(w)
}
