package invoicehdlr

import (
	"invoice-service/internal/domain"
	"invoice-service/pkg/server"
	"net/http"
)

type getTradesRequest struct {
	TradeStatus string
}

func (h *handler) GetTrades(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	/** Build request DTO */
	var req = getTradesRequest{
		TradeStatus: server.GetStringParam(r, "status", ""),
	}

	/** Service */
	res, err := h.service.FindTrades(ctx, &req.TradeStatus)
	if err != nil {
		server.InternalServerError(w, r, err)
		return
	}

	server.OK(w, r, buildTradeListResponse(res))
}

type tradesResponse struct {
	Investors []tradeBodyResponse `json:"trades"`
}

type tradeBodyResponse struct {
	ID           string  `json:"id"`
	InvoiceID    string  `json:"invoice_id"`
	InvestorsIDs []int   `json:"investors_ids"`
	TradeStatus  string  `json:"trade_status"`
	CreatedAt    string  `json:"created_at"`
	UpdatedAt    *string `json:"updated_at"`
}

func buildTradeListResponse(trades []domain.Trade) tradesResponse {
	var tradeListResponse tradesResponse
	for _, t := range trades {
		tradeListResponse.Investors = append(tradeListResponse.Investors, tradeBodyResponse{
			ID:           t.ID(),
			InvoiceID:    t.InvoiceID(),
			InvestorsIDs: t.InvestorsIDs(),
			TradeStatus:  t.TradeStatus(),
			CreatedAt:    t.CreatedAt(),
			UpdatedAt:    t.UpdatedAt(),
		})
	}

	return tradeListResponse
}
