package invoicehdlr

import (
	"encoding/json"
	"errors"
	"invoice-service/internal/application/invoicesrv"
	"invoice-service/internal/domain"
	"io/ioutil"
	"net/http"
)

type postInvoiceRequest struct {
	Body postInvoiceBodyRequest
}

type postInvoiceBodyRequest struct {
	Items    []itemBody `json:"items"`
	IssuerID int        `json:"issuer_id"`
}

type itemBody struct {
	ID          string  `json:"id"`
	Description string  `json:"description"`
	Price       float64 `json:"price"`
	Quantity    int     `json:"quantity"`
}

type invoiceResponse struct {
	ID           string     `json:"id"`
	DueDate      string     `json:"due_date"`
	AskingPrice  float64    `json:"asking_price"`
	Status       string     `json:"status"`
	Items        []itemBody `json:"items"`
	CreatedAt    string     `json:"created_at"`
	IssuerID     int        `json:"issuer_id"`
	InvestorsIDs []int      `json:"investors_ids"`
}

func newPostInvoiceRequest(r *http.Request) (*postInvoiceRequest, error) {
	var requestBody postInvoiceBodyRequest
	b, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &requestBody)
	if err != nil {
		return nil, err
	}

	var resp = &postInvoiceRequest{
		Body: requestBody,
	}

	return resp, nil
}

func (r *postInvoiceRequest) validate() error {
	if r.Body.IssuerID == -1 {
		return errors.New("issuer_id is required in body")
	}
	if r.Body.Items == nil {
		return errors.New("items is required in body")
	}
	return nil
}

func (p *postInvoiceBodyRequest) toInvoiceInput() *invoicesrv.InvoiceInput {
	var items []invoicesrv.ItemInput
	for _, i := range p.Items {
		items = append(items, invoicesrv.ItemInput{
			ID:          i.ID,
			Description: i.Description,
			Price:       i.Price,
			Quantity:    i.Quantity,
		})
	}

	return &invoicesrv.InvoiceInput{
		Items:    items,
		IssuerID: p.IssuerID,
	}
}

func buildInvoiceResponse(i *domain.Invoice, investorsIDs []int) *invoiceResponse {
	var items []itemBody
	for _, item := range i.Items() {
		items = append(items, itemBody{
			ID:          item.ID(),
			Description: item.Description(),
			Price:       item.Price(),
			Quantity:    item.Quantity(),
		})
	}

	return &invoiceResponse{
		ID:           i.ID(),
		DueDate:      i.DueDate(),
		AskingPrice:  i.AskingPrice(),
		Status:       i.Status(),
		Items:        items,
		CreatedAt:    i.CreatedAt(),
		IssuerID:     i.IssuerID(),
		InvestorsIDs: investorsIDs,
	}
}
