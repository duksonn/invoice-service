package invoicehdlr

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
)

type postBidRequest struct {
	Body postBidBodyRequest
}

type postBidBodyRequest struct {
	InvoiceID  string  `json:"invoice_id"`
	InvestorID int     `json:"investor_id"`
	BidAmount  float64 `json:"bid_amount"`
}

func newPostBidRequest(r *http.Request) (*postBidRequest, error) {
	var requestBody postBidBodyRequest
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &requestBody)
	if err != nil {
		return nil, err
	}

	var resp = &postBidRequest{
		Body: requestBody,
	}

	return resp, nil
}

func (r *postBidRequest) validate() error {
	if r.Body.InvestorID == -1 {
		return errors.New("investor_id is required in body")
	}
	if r.Body.InvoiceID == "" {
		return errors.New("invoice_id is required in body")
	}
	if r.Body.BidAmount == 0 {
		return errors.New("bid_amount is required in body")
	}
	return nil
}
