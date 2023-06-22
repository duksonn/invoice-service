package invoicehdlr

import (
	"encoding/json"
	"errors"
	"fmt"
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

func TestHandler_PutTrade(t *testing.T) {
	// Mock sunday day and patch
	mockSundayDate := time.Date(2022, 01, 30, 00, 00, 00, 00, time.UTC)
	utils.PatchNow(func() time.Time { return mockSundayDate })
	defer utils.RestoreNow()

	// Generate and patch mock ID
	utils.PatchGenerateUuid(func() string { return "id" })
	defer utils.RestoreGenerateUuid()

	type args struct {
		tradeID  *string
		approved *string
	}

	tests := []struct {
		name        string
		args        args
		mocks       func(dep *dependencies)
		errExpected *errResponse
	}{
		{
			name: "OK 1. Return nothing cause trade was approved correct",
			args: args{
				tradeID:  utils.PString("id"),
				approved: utils.PString("true"),
			},
			mocks: func(dep *dependencies) {
				var trade = domain.NewTrade("id", "invoiceID", []int{1}, "WAITING_APPROVAL", mockSundayDate.Format(time.RFC3339), nil)
				dep.databaseRepository.EXPECT().GetTradeByID(gomock.Any(), gomock.Any()).Return(trade, nil)

				dep.databaseRepository.EXPECT().SaveTrade(gomock.Any(), gomock.Any()).Return(nil, nil)

				invoice, _ := domain.NewInvoice(
					"id",
					mockSundayDate.AddDate(0, 0, 21).Format(time.RFC3339),
					35000,
					"CREATED",
					[]domain.Item{*domain.NewItem("id", "some desc", 500.0, 1)},
					mockSundayDate.Format(time.RFC3339),
					1,
				)
				dep.databaseRepository.EXPECT().GetInvoiceByID(gomock.Any(), gomock.Any()).Return(invoice, nil)

				var availableFunds = 1000.0
				var issuer = domain.NewIssuer(1, "some name", &availableFunds)
				dep.databaseRepository.EXPECT().GetIssuerByID(gomock.Any(), gomock.Any()).Return(issuer, nil)

				dep.databaseRepository.EXPECT().SaveIssuer(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			errExpected: nil,
		},
		{
			name: "OK 2. Return nothing cause trade was rejected correct",
			args: args{
				tradeID:  utils.PString("id"),
				approved: utils.PString("false"),
			},
			mocks: func(dep *dependencies) {
				var trade = domain.NewTrade("id", "invoiceID", []int{1}, "WAITING_APPROVAL", mockSundayDate.Format(time.RFC3339), nil)
				dep.databaseRepository.EXPECT().GetTradeByID(gomock.Any(), gomock.Any()).Return(trade, nil)

				dep.databaseRepository.EXPECT().SaveTrade(gomock.Any(), gomock.Any()).Return(nil, nil)

				invoice, _ := domain.NewInvoice(
					"id",
					mockSundayDate.AddDate(0, 0, 21).Format(time.RFC3339),
					35000,
					"CREATED",
					[]domain.Item{*domain.NewItem("id", "some desc", 500.0, 1)},
					mockSundayDate.Format(time.RFC3339),
					1,
				)
				dep.databaseRepository.EXPECT().GetInvoiceByID(gomock.Any(), gomock.Any()).Return(invoice, nil)

				var issuer = domain.NewIssuer(1, "some name", utils.PFloat(1000.0))
				dep.databaseRepository.EXPECT().GetIssuerByID(gomock.Any(), gomock.Any()).Return(issuer, nil)

				var investor = domain.NewInvestor(1, "some name", utils.PFloat(500.0))
				dep.databaseRepository.EXPECT().GetInvestorByID(gomock.Any(), gomock.Any()).Return(investor, nil)

				var bid = domain.NewBid("id", 1, "investor ID", 500, mockSundayDate.Format(time.RFC3339))
				dep.databaseRepository.EXPECT().GetBidsFromInvoiceAndInvestor(gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.Bid{*bid}, nil)

				dep.databaseRepository.EXPECT().SaveInvestor(gomock.Any(), gomock.Any()).Return(nil, nil)
			},
			errExpected: nil,
		},
		{
			name: "ERROR 1. Return error cause missing trade ID",
			args: args{
				tradeID:  nil,
				approved: utils.PString("false"),
			},
			mocks: func(dep *dependencies) {},
			errExpected: &errResponse{
				Messages: []string{"trade_id param is required"},
				Code:     "BAD_REQUEST",
			},
		},
		{
			name: "ERROR 2. Return error cause SaveInvestor failed",
			args: args{
				tradeID:  utils.PString("id"),
				approved: utils.PString("false"),
			},
			mocks: func(dep *dependencies) {
				var trade = domain.NewTrade("id", "invoiceID", []int{1}, "WAITING_APPROVAL", mockSundayDate.Format(time.RFC3339), nil)
				dep.databaseRepository.EXPECT().GetTradeByID(gomock.Any(), gomock.Any()).Return(trade, nil)

				dep.databaseRepository.EXPECT().SaveTrade(gomock.Any(), gomock.Any()).Return(nil, nil)

				invoice, _ := domain.NewInvoice(
					"id",
					mockSundayDate.AddDate(0, 0, 21).Format(time.RFC3339),
					35000,
					"CREATED",
					[]domain.Item{*domain.NewItem("id", "some desc", 500.0, 1)},
					mockSundayDate.Format(time.RFC3339),
					1,
				)
				dep.databaseRepository.EXPECT().GetInvoiceByID(gomock.Any(), gomock.Any()).Return(invoice, nil)

				var availableFunds = 1000.0
				var issuer = domain.NewIssuer(1, "some name", &availableFunds)
				dep.databaseRepository.EXPECT().GetIssuerByID(gomock.Any(), gomock.Any()).Return(issuer, nil)

				var investor = domain.NewInvestor(1, "some name", utils.PFloat(500.0))
				dep.databaseRepository.EXPECT().GetInvestorByID(gomock.Any(), gomock.Any()).Return(investor, nil)

				var bid = domain.NewBid("id", 1, "investor ID", 500, mockSundayDate.Format(time.RFC3339))
				dep.databaseRepository.EXPECT().GetBidsFromInvoiceAndInvestor(gomock.Any(), gomock.Any(), gomock.Any()).Return([]domain.Bid{*bid}, nil)

				dep.databaseRepository.EXPECT().SaveInvestor(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			errExpected: &errResponse{
				Messages: []string{"some error"},
				Code:     "INTERNAL_SERVER_ERROR",
			},
		},
		{
			name: "ERROR 3. Return error cause GetBidsFromInvoiceAndInvestor failed",
			args: args{
				tradeID:  utils.PString("id"),
				approved: utils.PString("false"),
			},
			mocks: func(dep *dependencies) {
				var trade = domain.NewTrade("id", "invoiceID", []int{1}, "WAITING_APPROVAL", mockSundayDate.Format(time.RFC3339), nil)
				dep.databaseRepository.EXPECT().GetTradeByID(gomock.Any(), gomock.Any()).Return(trade, nil)

				dep.databaseRepository.EXPECT().SaveTrade(gomock.Any(), gomock.Any()).Return(nil, nil)

				invoice, _ := domain.NewInvoice(
					"id",
					mockSundayDate.AddDate(0, 0, 21).Format(time.RFC3339),
					35000,
					"CREATED",
					[]domain.Item{*domain.NewItem("id", "some desc", 500.0, 1)},
					mockSundayDate.Format(time.RFC3339),
					1,
				)
				dep.databaseRepository.EXPECT().GetInvoiceByID(gomock.Any(), gomock.Any()).Return(invoice, nil)

				var availableFunds = 1000.0
				var issuer = domain.NewIssuer(1, "some name", &availableFunds)
				dep.databaseRepository.EXPECT().GetIssuerByID(gomock.Any(), gomock.Any()).Return(issuer, nil)

				var investor = domain.NewInvestor(1, "some name", utils.PFloat(500.0))
				dep.databaseRepository.EXPECT().GetInvestorByID(gomock.Any(), gomock.Any()).Return(investor, nil)

				dep.databaseRepository.EXPECT().GetBidsFromInvoiceAndInvestor(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			errExpected: &errResponse{
				Messages: []string{"some error"},
				Code:     "INTERNAL_SERVER_ERROR",
			},
		},
		{
			name: "ERROR 4. Return error cause GetInvestorByID failed",
			args: args{
				tradeID:  utils.PString("id"),
				approved: utils.PString("false"),
			},
			mocks: func(dep *dependencies) {
				var trade = domain.NewTrade("id", "invoiceID", []int{1}, "WAITING_APPROVAL", mockSundayDate.Format(time.RFC3339), nil)
				dep.databaseRepository.EXPECT().GetTradeByID(gomock.Any(), gomock.Any()).Return(trade, nil)

				dep.databaseRepository.EXPECT().SaveTrade(gomock.Any(), gomock.Any()).Return(nil, nil)

				invoice, _ := domain.NewInvoice(
					"id",
					mockSundayDate.AddDate(0, 0, 21).Format(time.RFC3339),
					35000,
					"CREATED",
					[]domain.Item{*domain.NewItem("id", "some desc", 500.0, 1)},
					mockSundayDate.Format(time.RFC3339),
					1,
				)
				dep.databaseRepository.EXPECT().GetInvoiceByID(gomock.Any(), gomock.Any()).Return(invoice, nil)

				var availableFunds = 1000.0
				var issuer = domain.NewIssuer(1, "some name", &availableFunds)
				dep.databaseRepository.EXPECT().GetIssuerByID(gomock.Any(), gomock.Any()).Return(issuer, nil)

				dep.databaseRepository.EXPECT().GetInvestorByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			errExpected: &errResponse{
				Messages: []string{"some error"},
				Code:     "INTERNAL_SERVER_ERROR",
			},
		},
		{
			name: "ERROR 5. Return error cause GetIssuerByID failed",
			args: args{
				tradeID:  utils.PString("id"),
				approved: utils.PString("false"),
			},
			mocks: func(dep *dependencies) {
				var trade = domain.NewTrade("id", "invoiceID", []int{1}, "WAITING_APPROVAL", mockSundayDate.Format(time.RFC3339), nil)
				dep.databaseRepository.EXPECT().GetTradeByID(gomock.Any(), gomock.Any()).Return(trade, nil)

				dep.databaseRepository.EXPECT().SaveTrade(gomock.Any(), gomock.Any()).Return(nil, nil)

				invoice, _ := domain.NewInvoice(
					"id",
					mockSundayDate.AddDate(0, 0, 21).Format(time.RFC3339),
					35000,
					"CREATED",
					[]domain.Item{*domain.NewItem("id", "some desc", 500.0, 1)},
					mockSundayDate.Format(time.RFC3339),
					1,
				)
				dep.databaseRepository.EXPECT().GetInvoiceByID(gomock.Any(), gomock.Any()).Return(invoice, nil)

				dep.databaseRepository.EXPECT().GetIssuerByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			errExpected: &errResponse{
				Messages: []string{"some error"},
				Code:     "INTERNAL_SERVER_ERROR",
			},
		},
		{
			name: "ERROR 6. Return error cause GetInvoiceByID failed",
			args: args{
				tradeID:  utils.PString("id"),
				approved: utils.PString("false"),
			},
			mocks: func(dep *dependencies) {
				var trade = domain.NewTrade("id", "invoiceID", []int{1}, "WAITING_APPROVAL", mockSundayDate.Format(time.RFC3339), nil)
				dep.databaseRepository.EXPECT().GetTradeByID(gomock.Any(), gomock.Any()).Return(trade, nil)

				dep.databaseRepository.EXPECT().SaveTrade(gomock.Any(), gomock.Any()).Return(nil, nil)

				dep.databaseRepository.EXPECT().GetInvoiceByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			errExpected: &errResponse{
				Messages: []string{"some error"},
				Code:     "INTERNAL_SERVER_ERROR",
			},
		},
		{
			name: "ERROR 7. Return error cause SaveTrade failed",
			args: args{
				tradeID:  utils.PString("id"),
				approved: utils.PString("false"),
			},
			mocks: func(dep *dependencies) {
				var trade = domain.NewTrade("id", "invoiceID", []int{1}, "WAITING_APPROVAL", mockSundayDate.Format(time.RFC3339), nil)
				dep.databaseRepository.EXPECT().GetTradeByID(gomock.Any(), gomock.Any()).Return(trade, nil)

				dep.databaseRepository.EXPECT().SaveTrade(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			errExpected: &errResponse{
				Messages: []string{"some error"},
				Code:     "INTERNAL_SERVER_ERROR",
			},
		},
		{
			name: "ERROR 8. Return error cause GetTradeByID failed",
			args: args{
				tradeID:  utils.PString("id"),
				approved: utils.PString("false"),
			},
			mocks: func(dep *dependencies) {
				dep.databaseRepository.EXPECT().GetTradeByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			errExpected: &errResponse{
				Messages: []string{"some error"},
				Code:     "INTERNAL_SERVER_ERROR",
			},
		},
		{
			name: "ERROR 9. Return error cause SaveIssuer failed",
			args: args{
				tradeID:  utils.PString("id"),
				approved: utils.PString("true"),
			},
			mocks: func(dep *dependencies) {
				var trade = domain.NewTrade("id", "invoiceID", []int{1}, "WAITING_APPROVAL", mockSundayDate.Format(time.RFC3339), nil)
				dep.databaseRepository.EXPECT().GetTradeByID(gomock.Any(), gomock.Any()).Return(trade, nil)

				dep.databaseRepository.EXPECT().SaveTrade(gomock.Any(), gomock.Any()).Return(nil, nil)

				invoice, _ := domain.NewInvoice(
					"id",
					mockSundayDate.AddDate(0, 0, 21).Format(time.RFC3339),
					35000,
					"CREATED",
					[]domain.Item{*domain.NewItem("id", "some desc", 500.0, 1)},
					mockSundayDate.Format(time.RFC3339),
					1,
				)
				dep.databaseRepository.EXPECT().GetInvoiceByID(gomock.Any(), gomock.Any()).Return(invoice, nil)

				var availableFunds = 1000.0
				var issuer = domain.NewIssuer(1, "some name", &availableFunds)
				dep.databaseRepository.EXPECT().GetIssuerByID(gomock.Any(), gomock.Any()).Return(issuer, nil)

				dep.databaseRepository.EXPECT().SaveIssuer(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
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

			var path string
			var route string
			if tt.args.tradeID != nil && *tt.args.tradeID != "" {
				path = fmt.Sprintf("/v1/trade/%v", tt.args.tradeID)
				route = "/v1/trade/{trade_id}"
			} else {
				path = "/v1/trade"
				route = "/v1/trade"
			}

			// Create reader with URL and params
			var method = http.MethodPut
			var reader = httptest.NewRequest(method, path, nil)
			var q = reader.URL.Query()
			if tt.args.approved != nil && *tt.args.approved != "" {
				q.Add("approved", *tt.args.approved)
			}
			reader.URL.RawQuery = q.Encode()

			// Create writer
			var writer = httptest.NewRecorder()

			// Build custom server
			var router = mux.NewRouter()
			router.HandleFunc(route, handler.PutTrade).Methods(method)
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
