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

func TestHandler_PostInvoice(t *testing.T) {
	// Mock sunday day and patch
	mockSundayDate := time.Date(2022, 01, 30, 00, 00, 00, 00, time.UTC)
	utils.PatchNow(func() time.Time { return mockSundayDate })
	defer utils.RestoreNow()

	// Generate and patch mock ID
	utils.PatchGenerateUuid(func() string { return "id" })
	defer utils.RestoreGenerateUuid()

	tests := []struct {
		name        string
		args        map[string]interface{}
		mocks       func(dep *dependencies)
		expected    map[string]interface{}
		errExpected *errResponse
	}{
		{
			name: "OK 1. Return invoice",
			args: map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"id":          "item ID",
						"description": "some desc",
						"price":       500.0,
						"quantity":    1,
					},
				},
				"issuer_id": 1,
			},
			mocks: func(dep *dependencies) {
				invoice, _ := domain.NewInvoice(
					"id",
					mockSundayDate.AddDate(0, 0, 21).Format(time.RFC3339),
					500.0,
					"CREATED",
					[]domain.Item{*domain.NewItem("item ID", "some desc", 500.0, 1)},
					mockSundayDate.Format(time.RFC3339),
					1,
				)
				dep.databaseRepository.EXPECT().SaveInvoice(gomock.Any(), gomock.Any()).Return(invoice, nil)
			},
			expected: map[string]interface{}{
				"id":           "id",
				"due_date":     mockSundayDate.AddDate(0, 0, 21),
				"asking_price": 500.0,
				"status":       "CREATED",
				"items": []map[string]interface{}{
					{
						"id":          "item ID",
						"description": "some desc",
						"price":       500,
						"quantity":    1,
					},
				},
				"created_at": mockSundayDate,
				"issuer_id":  1,
			},
			errExpected: nil,
		},
		{
			name: "ERROR 1. Return error cause items is missing",
			args: map[string]interface{}{
				"issuer_id": 1,
			},
			mocks:    func(dep *dependencies) {},
			expected: nil,
			errExpected: &errResponse{
				Messages: []string{"items is required in body"},
				Code:     "BAD_REQUEST",
			},
		},
		{
			name: "ERROR 2. Return error cause issuer ID is missing",
			args: map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"id":          "item ID",
						"description": "some desc",
						"price":       500.0,
						"quantity":    1,
					},
				},
				"issuer_id": -1,
			},
			mocks:    func(dep *dependencies) {},
			expected: nil,
			errExpected: &errResponse{
				Messages: []string{"issuer_id is required in body"},
				Code:     "BAD_REQUEST",
			},
		},
		{
			name: "ERROR 3. Return error cause issuer ID is missing",
			args: map[string]interface{}{
				"items": []map[string]interface{}{
					{
						"id":          "item ID",
						"description": "some desc",
						"price":       500.0,
						"quantity":    1,
					},
				},
				"issuer_id": 1,
			},
			mocks:    func(dep *dependencies) {
				dep.databaseRepository.EXPECT().SaveInvoice(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			expected: nil,
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
			var path = "/v1/invoice"
			var method = http.MethodPost
			req, _ := json.Marshal(tt.args)
			var reader = httptest.NewRequest(method, path, bytes.NewReader(req))
			var q = reader.URL.Query()
			reader.URL.RawQuery = q.Encode()

			// Create writer
			var writer = httptest.NewRecorder()

			// Build custom server
			var router = mux.NewRouter()
			router.HandleFunc(path, handler.PostInvoice).Methods(method)
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
