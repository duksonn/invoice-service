package invoicehdlr

import (
	"encoding/json"
	"errors"
	"github.com/golang/mock/gomock"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"
	"invoice-service/internal/domain"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandler_GetInvestors(t *testing.T) {
	tests := []struct {
		name        string
		mocks       func(dep *dependencies)
		expected    map[string]interface{}
		errExpected *errResponse
	}{
		{
			name: "OK 1. Return investors",
			mocks: func(dep *dependencies) {
				var availableFunds = 1000.0
				var investor = domain.NewInvestor(1, "some name", &availableFunds)
				dep.databaseRepository.EXPECT().FindInvestors(gomock.Any()).Return([]domain.Investor{*investor}, nil)
			},
			expected: map[string]interface{}{
				"investors": []map[string]interface{}{
					{
						"id":              1,
						"name":            "some name",
						"available_funds": 1000.0,
					},
				},
			},
			errExpected: nil,
		},
		{
			name: "ERROR 1. Failed getting investors from DB",
			mocks: func(dep *dependencies) {
				dep.databaseRepository.EXPECT().FindInvestors(gomock.Any()).Return(nil, errors.New("some error"))
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
			var path = "/v1/investor"
			var method = http.MethodGet
			var reader = httptest.NewRequest(method, path, nil)
			var q = reader.URL.Query()
			reader.URL.RawQuery = q.Encode()

			// Create writer
			var writer = httptest.NewRecorder()

			// Build custom server
			var router = mux.NewRouter()
			router.HandleFunc(path, handler.GetInvestors).Methods(method)
			router.ServeHTTP(writer, reader)

			// Asserts
			assert.NotNil(t, handler.service)
			if tt.expected != nil {
				expectedRespByte, _ := json.Marshal(tt.expected)
				var expectedResp investorsResponse
				_ = json.Unmarshal(expectedRespByte, &expectedResp)

				var resp investorsResponse
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
