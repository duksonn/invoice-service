package invoicehdlr

import (
	"invoice-service/pkg/server"
	"net/http"
	"strconv"
)

type putTradeRequest struct {
	TradeID  string
	Approved string
}

func (h *handler) PutTrade(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	/** Build request DTO */
	var req = putTradeRequest{
		TradeID:  server.GetStringFromPath(r, "trade_id", ""),
		Approved: server.GetStringParam(r, "approved", ""),
	}

	/** Validations */
	if req.TradeID == "" {
		server.BadRequest(w, r, "BAD_REQUEST", "trade_id param is required")
		return
	}

	isApproved, err := strconv.ParseBool(req.Approved)
	if err != nil {
		server.InternalServerError(w, r, err)
		return
	}

	/** Service */
	err = h.service.ApproveOrRejectTrade(ctx, req.TradeID, isApproved)
	if err != nil {
		server.InternalServerError(w, r, err)
		return
	}

	server.OKNoContent(w)
}
