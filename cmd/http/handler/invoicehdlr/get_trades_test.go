package invoicehdlr

import (
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

func TestHandler_GetTrades(t *testing.T) {
	// Mock sunday day and patch
	mockSundayDate := time.Date(2022, 01, 30, 00, 00, 00, 00, time.UTC)
	utils.PatchNow(func() time.Time { return mockSundayDate })
	defer utils.RestoreNow()

	tests := []struct {
		name        string
		args        string
		mocks       func(dep *dependencies)
		expected    map[string]interface{}
		errExpected *errResponse
	}{
		{
			name: "OK 1. Return trades",
			args: "WAITING_APPROVAL",
			mocks: func(dep *dependencies) {
				var trade = domain.NewTrade(
					"id",
					"invoiceID",
					[]int{1},
					"WAITING_APPROVAL",
					mockSundayDate.Format(time.RFC3339),
					nil,
				)
				dep.databaseRepository.EXPECT().FindTrades(gomock.Any(), gomock.Any()).Return([]domain.Trade{*trade}, nil)
			},
			expected: map[string]interface{}{
				"trades": []map[string]interface{}{
					{
						"id":            "id",
						"invoice_id":    "invoiceID",
						"investors_ids": []int{1},
						"trade_status":  "WAITING_APPROVAL",
						"created_at":    mockSundayDate,
					},
				},
			},
			errExpected: nil,
		},
		{
			name: "OK 2. Return empty cause no found any trade with status ACCEPTED",
			args: "ACCEPTED",
			mocks: func(dep *dependencies) {
				dep.databaseRepository.EXPECT().FindTrades(gomock.Any(), gomock.Any()).Return([]domain.Trade{}, nil)
			},
			expected: map[string]interface{}{
				"trades": nil,
			},
			errExpected: nil,
		},
		{
			name: "ERROR 1. Failed getting trades from DB",
			args: "WAITING_APPROVAL",
			mocks: func(dep *dependencies) {
				dep.databaseRepository.EXPECT().FindTrades(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			expected: nil,
			errExpected: &errResponse{
				Messages: []string{"some error"},
				Code:     "INTERNAL_SERVER_ERROR",
			},
		},
		{
			name:     "ERROR 2. Failed getting trades cause status is not valid",
			args:     "WAITING_APPROVALLlLllL",
			mocks:    func(dep *dependencies) {},
			expected: nil,
			errExpected: &errResponse{
				Messages: []string{"invalid trade status"},
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
			var path = "/v1/trade"
			var method = http.MethodGet
			var reader = httptest.NewRequest(method, path, nil)
			var q = reader.URL.Query()
			q.Add("status", tt.args)
			reader.URL.RawQuery = q.Encode()

			// Create writer
			var writer = httptest.NewRecorder()

			// Build custom server
			var router = mux.NewRouter()
			router.HandleFunc(path, handler.GetTrades).Methods(method)
			router.ServeHTTP(writer, reader)

			// Asserts
			assert.NotNil(t, handler.service)
			if tt.expected != nil {
				expectedRespByte, _ := json.Marshal(tt.expected)
				var expectedResp tradesResponse
				_ = json.Unmarshal(expectedRespByte, &expectedResp)

				var resp tradesResponse
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
