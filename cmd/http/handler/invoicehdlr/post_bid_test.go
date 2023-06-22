package invoicehdlr

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"invoice-service/internal/domain"
	"invoice-service/pkg/utils"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestHandler_PostBid(t *testing.T) {
	// Mock sunday day and patch
	mockSundayDate := time.Date(2022, 01, 30, 00, 00, 00, 00, time.UTC)
	utils.PatchNow(func() time.Time { return mockSundayDate })
	defer utils.RestoreNow()

	tests := []struct {
		name        string
		args        map[string]interface{}
		mocks       func(dep *dependencies)
		errExpected *errResponse
	}{
		{
			name: "OK 1. Return nothing cause bid had been placed",
			args: map[string]interface{}{
				"invoice_id":  "id",
				"investor_id": 1,
				"bid_amount":  500,
			},
			mocks: func(dep *dependencies) {
				var investor = domain.NewInvestor(1, "some name", utils.PFloat(500.0))
				dep.databaseRepository.EXPECT().GetInvestorByID(gomock.Any(), gomock.Any()).Return(investor, nil)

				dep.databaseRepository.EXPECT().SaveInvestor(gomock.Any(), gomock.Any()).Return(nil, nil)

				var bid = domain.NewBid("id", 1, "investor ID", 500, mockSundayDate.Format(time.RFC3339))
				dep.databaseRepository.EXPECT().SaveBid(gomock.Any(), gomock.Any()).Return(bid, nil)

				dep.databaseRepository.EXPECT().GetTotalBidsAmount(gomock.Any(), gomock.Any()).Return(utils.PFloat(600.0), nil)

				invoice, _ := domain.NewInvoice(
					"id",
					mockSundayDate.AddDate(0, 0, 21).Format(time.RFC3339),
					500.0,
					"CREATED",
					[]domain.Item{*domain.NewItem("id", "some desc", 500.0, 1)},
					mockSundayDate.Format(time.RFC3339),
					1,
				)
				dep.databaseRepository.EXPECT().GetInvoiceByID(gomock.Any(), gomock.Any()).Return(invoice, nil)

				dep.databaseRepository.EXPECT().SaveInvoice(gomock.Any(), gomock.Any()).Return(nil, nil)

				dep.databaseRepository.EXPECT().GetBidsFromInvoiceAndInvestor(gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.Bid{*bid}, nil)

				dep.databaseRepository.EXPECT().SaveTrade(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			errExpected: nil,
		},
		{
			name: "OK 2. Return nothing cause bid had been placed but not purchase the invoice",
			args: map[string]interface{}{
				"invoice_id":  "id",
				"investor_id": 1,
				"bid_amount":  500,
			},
			mocks: func(dep *dependencies) {
				var investor = domain.NewInvestor(1, "some name", utils.PFloat(500.0))
				dep.databaseRepository.EXPECT().GetInvestorByID(gomock.Any(), gomock.Any()).Return(investor, nil)

				dep.databaseRepository.EXPECT().SaveInvestor(gomock.Any(), gomock.Any()).Return(nil, nil)

				var bid = domain.NewBid("id", 1, "investor ID", 500, mockSundayDate.Format(time.RFC3339))
				dep.databaseRepository.EXPECT().SaveBid(gomock.Any(), gomock.Any()).Return(bid, nil)

				dep.databaseRepository.EXPECT().GetTotalBidsAmount(gomock.Any(), gomock.Any()).Return(utils.PFloat(600.0), nil)

				invoice, _ := domain.NewInvoice(
					"id",
					mockSundayDate.AddDate(0, 0, 21).Format(time.RFC3339),
					3500.0,
					"CREATED",
					[]domain.Item{*domain.NewItem("id", "some desc", 500.0, 1)},
					mockSundayDate.Format(time.RFC3339),
					1,
				)
				dep.databaseRepository.EXPECT().GetInvoiceByID(gomock.Any(), gomock.Any()).Return(invoice, nil)
			},
			errExpected: nil,
		},
		{
			name: "ERROR 1. Return error cause investor ID must be valid",
			args: map[string]interface{}{
				"invoice_id":  "id",
				"investor_id": -1,
				"bid_amount":  500,
			},
			mocks: func(dep *dependencies) {},
			errExpected: &errResponse{
				Messages: []string{"investor_id is required in body"},
				Code:     "BAD_REQUEST",
			},
		},
		{
			name: `ERROR 2. Return error cause invoice ID must be valid`,
			args: map[string]interface{}{
				"invoice_id":  "",
				"investor_id": 1,
				"bid_amount":  500,
			},
			mocks: func(dep *dependencies) {},
			errExpected: &errResponse{
				Messages: []string{"invoice_id is required in body"},
				Code:     "BAD_REQUEST",
			},
		},
		{
			name: `ERROR 3. Return error cause bid amount must be valid`,
			args: map[string]interface{}{
				"invoice_id":  "id",
				"investor_id": 1,
				"bid_amount":  0,
			},
			mocks: func(dep *dependencies) {},
			errExpected: &errResponse{
				Messages: []string{"bid_amount is required in body"},
				Code:     "BAD_REQUEST",
			},
		},
		{
			name: "ERROR 4. Return error cause GetInvestorByID failed",
			args: map[string]interface{}{
				"invoice_id":  "id",
				"investor_id": 1,
				"bid_amount":  500,
			},
			mocks: func(dep *dependencies) {
				dep.databaseRepository.EXPECT().GetInvestorByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			errExpected: &errResponse{
				Messages: []string{"some error"},
				Code:     "INTERNAL_SERVER_ERROR",
			},
		},
		{
			name: "ERROR 5. Return error cause SaveInvestor failed",
			args: map[string]interface{}{
				"invoice_id":  "id",
				"investor_id": 1,
				"bid_amount":  500,
			},
			mocks: func(dep *dependencies) {
				var investor = domain.NewInvestor(1, "some name", utils.PFloat(500.0))
				dep.databaseRepository.EXPECT().GetInvestorByID(gomock.Any(), gomock.Any()).Return(investor, nil)

				dep.databaseRepository.EXPECT().SaveInvestor(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			errExpected: &errResponse{
				Messages: []string{"some error"},
				Code:     "INTERNAL_SERVER_ERROR",
			},
		},
		{
			name: "ERROR 6. Return error cause SaveBid failed",
			args: map[string]interface{}{
				"invoice_id":  "id",
				"investor_id": 1,
				"bid_amount":  500,
			},
			mocks: func(dep *dependencies) {
				var investor = domain.NewInvestor(1, "some name", utils.PFloat(500.0))
				dep.databaseRepository.EXPECT().GetInvestorByID(gomock.Any(), gomock.Any()).Return(investor, nil)

				dep.databaseRepository.EXPECT().SaveInvestor(gomock.Any(), gomock.Any()).Return(nil, nil)

				dep.databaseRepository.EXPECT().SaveBid(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			errExpected: &errResponse{
				Messages: []string{"some error"},
				Code:     "INTERNAL_SERVER_ERROR",
			},
		},
		{
			name: "ERROR 7. Return error cause GetTotalBidsAmount failed",
			args: map[string]interface{}{
				"invoice_id":  "id",
				"investor_id": 1,
				"bid_amount":  500,
			},
			mocks: func(dep *dependencies) {
				var investor = domain.NewInvestor(1, "some name", utils.PFloat(500.0))
				dep.databaseRepository.EXPECT().GetInvestorByID(gomock.Any(), gomock.Any()).Return(investor, nil)

				dep.databaseRepository.EXPECT().SaveInvestor(gomock.Any(), gomock.Any()).Return(nil, nil)

				var bid = domain.NewBid("id", 1, "investor ID", 500, mockSundayDate.Format(time.RFC3339))
				dep.databaseRepository.EXPECT().SaveBid(gomock.Any(), gomock.Any()).Return(bid, nil)

				dep.databaseRepository.EXPECT().GetTotalBidsAmount(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			errExpected: &errResponse{
				Messages: []string{"some error"},
				Code:     "INTERNAL_SERVER_ERROR",
			},
		},
		{
			name: "ERROR 8. Return error cause GetInvoiceByID failed",
			args: map[string]interface{}{
				"invoice_id":  "id",
				"investor_id": 1,
				"bid_amount":  500,
			},
			mocks: func(dep *dependencies) {
				var investor = domain.NewInvestor(1, "some name", utils.PFloat(500.0))
				dep.databaseRepository.EXPECT().GetInvestorByID(gomock.Any(), gomock.Any()).Return(investor, nil)

				dep.databaseRepository.EXPECT().SaveInvestor(gomock.Any(), gomock.Any()).Return(nil, nil)

				var bid = domain.NewBid("id", 1, "investor ID", 500, mockSundayDate.Format(time.RFC3339))
				dep.databaseRepository.EXPECT().SaveBid(gomock.Any(), gomock.Any()).Return(bid, nil)

				dep.databaseRepository.EXPECT().GetTotalBidsAmount(gomock.Any(), gomock.Any()).Return(utils.PFloat(600.0), nil)

				dep.databaseRepository.EXPECT().GetInvoiceByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			errExpected: &errResponse{
				Messages: []string{"some error"},
				Code:     "INTERNAL_SERVER_ERROR",
			},
		},
		{
			name: "ERROR 9. Return error cause SaveInvoice failed",
			args: map[string]interface{}{
				"invoice_id":  "id",
				"investor_id": 1,
				"bid_amount":  500,
			},
			mocks: func(dep *dependencies) {
				var investor = domain.NewInvestor(1, "some name", utils.PFloat(500.0))
				dep.databaseRepository.EXPECT().GetInvestorByID(gomock.Any(), gomock.Any()).Return(investor, nil)

				dep.databaseRepository.EXPECT().SaveInvestor(gomock.Any(), gomock.Any()).Return(nil, nil)

				var bid = domain.NewBid("id", 1, "investor ID", 500, mockSundayDate.Format(time.RFC3339))
				dep.databaseRepository.EXPECT().SaveBid(gomock.Any(), gomock.Any()).Return(bid, nil)

				dep.databaseRepository.EXPECT().GetTotalBidsAmount(gomock.Any(), gomock.Any()).Return(utils.PFloat(600.0), nil)

				invoice, _ := domain.NewInvoice(
					"id",
					mockSundayDate.AddDate(0, 0, 21).Format(time.RFC3339),
					500.0,
					"CREATED",
					[]domain.Item{*domain.NewItem("id", "some desc", 500.0, 1)},
					mockSundayDate.Format(time.RFC3339),
					1,
				)
				dep.databaseRepository.EXPECT().GetInvoiceByID(gomock.Any(), gomock.Any()).Return(invoice, nil)

				dep.databaseRepository.EXPECT().SaveInvoice(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			errExpected: &errResponse{
				Messages: []string{"some error"},
				Code:     "INTERNAL_SERVER_ERROR",
			},
		},
		{
			name: "ERROR 10. Return error cause GetBidsFromInvoiceAndInvestor failed",
			args: map[string]interface{}{
				"invoice_id":  "id",
				"investor_id": 1,
				"bid_amount":  500,
			},
			mocks: func(dep *dependencies) {
				var investor = domain.NewInvestor(1, "some name", utils.PFloat(500.0))
				dep.databaseRepository.EXPECT().GetInvestorByID(gomock.Any(), gomock.Any()).Return(investor, nil)

				dep.databaseRepository.EXPECT().SaveInvestor(gomock.Any(), gomock.Any()).Return(nil, nil)

				var bid = domain.NewBid("id", 1, "investor ID", 500, mockSundayDate.Format(time.RFC3339))
				dep.databaseRepository.EXPECT().SaveBid(gomock.Any(), gomock.Any()).Return(bid, nil)

				dep.databaseRepository.EXPECT().GetTotalBidsAmount(gomock.Any(), gomock.Any()).Return(utils.PFloat(600.0), nil)

				invoice, _ := domain.NewInvoice(
					"id",
					mockSundayDate.AddDate(0, 0, 21).Format(time.RFC3339),
					500.0,
					"CREATED",
					[]domain.Item{*domain.NewItem("id", "some desc", 500.0, 1)},
					mockSundayDate.Format(time.RFC3339),
					1,
				)
				dep.databaseRepository.EXPECT().GetInvoiceByID(gomock.Any(), gomock.Any()).Return(invoice, nil)

				dep.databaseRepository.EXPECT().SaveInvoice(gomock.Any(), gomock.Any()).Return(nil, nil)

				dep.databaseRepository.EXPECT().GetBidsFromInvoiceAndInvestor(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			errExpected: &errResponse{
				Messages: []string{"some error"},
				Code:     "INTERNAL_SERVER_ERROR",
			},
		},
		{
			name: "ERROR 11. Return error cause SaveTrade failed",
			args: map[string]interface{}{
				"invoice_id":  "id",
				"investor_id": 1,
				"bid_amount":  500,
			},
			mocks: func(dep *dependencies) {
				var investor = domain.NewInvestor(1, "some name", utils.PFloat(500.0))
				dep.databaseRepository.EXPECT().GetInvestorByID(gomock.Any(), gomock.Any()).Return(investor, nil)

				dep.databaseRepository.EXPECT().SaveInvestor(gomock.Any(), gomock.Any()).Return(nil, nil)

				var bid = domain.NewBid("id", 1, "investor ID", 500, mockSundayDate.Format(time.RFC3339))
				dep.databaseRepository.EXPECT().SaveBid(gomock.Any(), gomock.Any()).Return(bid, nil)

				dep.databaseRepository.EXPECT().GetTotalBidsAmount(gomock.Any(), gomock.Any()).Return(utils.PFloat(600.0), nil)

				invoice, _ := domain.NewInvoice(
					"id",
					mockSundayDate.AddDate(0, 0, 21).Format(time.RFC3339),
					500.0,
					"CREATED",
					[]domain.Item{*domain.NewItem("id", "some desc", 500.0, 1)},
					mockSundayDate.Format(time.RFC3339),
					1,
				)
				dep.databaseRepository.EXPECT().GetInvoiceByID(gomock.Any(), gomock.Any()).Return(invoice, nil)

				dep.databaseRepository.EXPECT().SaveInvoice(gomock.Any(), gomock.Any()).Return(nil, nil)

				dep.databaseRepository.EXPECT().GetBidsFromInvoiceAndInvestor(gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.Bid{*bid}, nil)

				dep.databaseRepository.EXPECT().SaveTrade(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			errExpected: &errResponse{
				Messages: []string{"some error"},
				Code:     "INTERNAL_SERVER_ERROR",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var dependencies = makeDependencies(t)
			var handler = buildHandler(dependencies)
			tt.mocks(dependencies)

			// Create reader with URL and params
			var path = "/v1/bid/place"
			var method = http.MethodPost
			req, _ := json.Marshal(tt.args)
			var reader = httptest.NewRequest(method, path, bytes.NewReader(req))
			var q = reader.URL.Query()
			reader.URL.RawQuery = q.Encode()

			// Create writer
			var writer = httptest.NewRecorder()

			// Build custom server
			var router = mux.NewRouter()
			router.HandleFunc(path, handler.PostBid).Methods(method)
			router.ServeHTTP(writer, reader)

			// Asserts
			assert.NotNil(t, handler.service)
			if tt.errExpected != nil {
				var resp errResponse
				_ = json.Unmarshal(writer.Body.Bytes(), &resp)

				assert.Equal(t, *tt.errExpected, resp)
			}
		})
	}
}
