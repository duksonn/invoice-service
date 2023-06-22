package invoicehdlr

import (
	"errors"
	"invoice-service/internal/domain"
	"invoice-service/pkg/server"
	"net/http"
)

type getIssuerRequest struct {
	IssuerID string
}

func (h *handler) GetIssuer(w http.ResponseWriter, r *http.Request) {
	var ctx = r.Context()

	/** Build request DTO */
	var req = getIssuerRequest{
		IssuerID: server.GetStringFromPath(r, "issuer_id", ""),
	}

	/** Validate request */
	var err = req.validate()
	if err != nil {
		server.BadRequest(w, r, "BAD_REQUEST", err.Error())
		return
	}

	/** Service */
	res, err := h.service.GetIssuer(ctx, req.IssuerID)
	if err != nil {
		server.InternalServerError(w, r, err)
		return
	}

	server.OK(w, r, buildIssuerResponse(res))
}

func (r *getIssuerRequest) validate() error {
	if r.IssuerID == "" {
		return errors.New("issuer_id param is required")
	}
	return nil
}

type issuerResponse struct {
	ID             int      `json:"id"`
	CompanyName    string   `json:"company_name"`
	AvailableFunds *float64 `json:"available_funds"`
}

func buildIssuerResponse(i *domain.Issuer) issuerResponse {
	return issuerResponse{
		ID:             i.ID(),
		CompanyName:    i.CompanyName(),
		AvailableFunds: i.AvailableFunds(),
	}
}
