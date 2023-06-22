package invoicehdlr

import (
	"invoice-service/internal/domain"
	"invoice-service/pkg/server"
	"net/http"
)

func (h *handler) GetInvestors(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	/** Service */
	res, err := h.service.FindInvestors(ctx)
	if err != nil {
		server.InternalServerError(w, r, err)
		return
	}

	server.OK(w, r, buildInvestorListResponse(res))
}

type investorsResponse struct {
	Investors []investorBodyResponse `json:"investors"`
}

type investorBodyResponse struct {
	ID             int      `json:"id"`
	Name           string   `json:"name"`
	AvailableFunds *float64 `json:"available_funds"`
}

func buildInvestorListResponse(investors []domain.Investor) investorsResponse {
	var investorListResponse investorsResponse
	for _, i := range investors {
		investorListResponse.Investors = append(investorListResponse.Investors, investorBodyResponse{
			ID:             i.ID(),
			Name:           i.Name(),
			AvailableFunds: i.AvailableFunds(),
		})
	}

	return investorListResponse
}
