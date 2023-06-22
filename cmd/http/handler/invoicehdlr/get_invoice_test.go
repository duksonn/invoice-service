package invoicehdlr

import (
	"database/sql"
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

func TestHandler_GetInvoice(t *testing.T) {
	// Mock sunday day and patch
	mockSundayDate := time.Date(2022, 01, 30, 00, 00, 00, 00, time.UTC)
	utils.PatchNow(func() time.Time { return mockSundayDate })
	defer utils.RestoreNow()

	tests := []struct {
		name        string
		args        *string
		mocks       func(dep *dependencies)
		expected    map[string]interface{}
		errExpected *errResponse
	}{
		{
			name: "OK 1. Return invoice",
			args: utils.PString("id"),
			mocks: func(dep *dependencies) {
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

				dep.databaseRepository.EXPECT().GetTradeByInvoice(gomock.Any(), gomock.Any()).Return(nil, sql.ErrNoRows)
			},
			expected: map[string]interface{}{
				"id":           "id",
				"due_date":     mockSundayDate.AddDate(0, 0, 21),
				"asking_price": 35000,
				"status":       "CREATED",
				"items": []map[string]interface{}{
					{
						"id":          "id",
						"description": "some desc",
						"price":       500.0,
						"quantity":    1,
					},
				},
				"created_at": mockSundayDate,
				"issuer_id":  1,
			},
			errExpected: nil,
		},
		{
			name: "OK 1. Return a purchased invoice",
			args: utils.PString("id"),
			mocks: func(dep *dependencies) {
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

				var trade = domain.NewTrade("id", "invoiceID", []int{1}, "APPROVED", mockSundayDate.Format(time.RFC3339), nil)
				dep.databaseRepository.EXPECT().GetTradeByInvoice(gomock.Any(), gomock.Any()).Return(trade, nil)
			},
			expected: map[string]interface{}{
				"id":           "id",
				"due_date":     mockSundayDate.AddDate(0, 0, 21),
				"asking_price": 35000,
				"status":       "CREATED",
				"items": []map[string]interface{}{
					{
						"id":          "id",
						"description": "some desc",
						"price":       500.0,
						"quantity":    1,
					},
				},
				"created_at":    mockSundayDate,
				"issuer_id":     1,
				"investors_ids": []int{1},
			},
			errExpected: nil,
		},
		{
			name: "ERROR 1. Failed getting invoice from DB",
			args: utils.PString("id"),
			mocks: func(dep *dependencies) {
				dep.databaseRepository.EXPECT().GetInvoiceByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			expected: nil,
			errExpected: &errResponse{
				Messages: []string{"some error"},
				Code:     "INTERNAL_SERVER_ERROR",
			},
		},
		{
			name:     "ERROR 2. Failed cause invoice ID is missing",
			args:     nil,
			mocks:    func(dep *dependencies) {},
			expected: nil,
			errExpected: &errResponse{
				Messages: []string{"invoice_id param is required"},
				Code:     "BAD_REQUEST",
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
			if tt.args != nil && *tt.args != "" {
				path = fmt.Sprintf("/v1/invoice/%v", *tt.args)
				route = "/v1/invoice/{invoice_id}"
			} else {
				path = "/v1/invoice"
				route = "/v1/invoice"
			}

			// Create reader with URL and params
			var method = http.MethodGet
			var reader = httptest.NewRequest(method, path, nil)
			var q = reader.URL.Query()
			reader.URL.RawQuery = q.Encode()

			// Create writer
			var writer = httptest.NewRecorder()

			// Build custom server
			var router = mux.NewRouter()
			router.HandleFunc(route, handler.GetInvoice).Methods(method)
			router.ServeHTTP(writer, reader)

			// Asserts
			assert.NotNil(t, handler.service)
			if tt.expected != nil {
				expectedRespByte, _ := json.Marshal(tt.expected)
				var expectedResp invoiceResponse
				_ = json.Unmarshal(expectedRespByte, &expectedResp)

				var resp invoiceResponse
				_ = json.Unmarshal(writer.Body.Bytes(), &resp)

				assert.Equal(t, expectedResp, resp)
			}
			if tt.errExpected != nil {
				var resp errResponse
				_ = json.Unmarshal(writer.Body.Bytes(), &resp)

				assert.Equal(t, *tt.errExpected, resp)
			}
		})
	}
}
