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
)

func TestHandler_GetIssuer(t *testing.T) {
	tests := []struct {
		name        string
		args        *string
		mocks       func(dep *dependencies)
		expected    map[string]interface{}
		errExpected *errResponse
	}{
		{
			name: "OK 1. Return issuer",
			args: utils.PString("1"),
			mocks: func(dep *dependencies) {
				var issuer = domain.NewIssuer(1, "some name", utils.PFloat(1000.0))
				dep.databaseRepository.EXPECT().GetIssuerByID(gomock.Any(), gomock.Any()).Return(issuer, nil)
			},
			expected: map[string]interface{}{
				"id":              1,
				"company_name":    "some name",
				"available_funds": 1000,
			},
			errExpected: nil,
		},
		{
			name: "ERROR 1. Failed getting issuer from DB",
			args: utils.PString("1"),
			mocks: func(dep *dependencies) {
				dep.databaseRepository.EXPECT().GetIssuerByID(gomock.Any(), gomock.Any()).Return(nil, errors.New("some error"))
			},
			expected: nil,
			errExpected: &errResponse{
				Messages: []string{"some error"},
				Code:     "INTERNAL_SERVER_ERROR",
			},
		},
		{
			name:     "ERROR 2. Failed cause issuer ID is missing",
			args:     nil,
			mocks:    func(dep *dependencies) {},
			expected: nil,
			errExpected: &errResponse{
				Messages: []string{"issuer_id param is required"},
				Code:     "BAD_REQUEST",
			},
		},
		{
			name: "ERROR 3. Failed getting issuer from DB",
			args: utils.PString("id"),
			mocks: func(dep *dependencies) {},
			expected: nil,
			errExpected: &errResponse{
				Messages: []string{"strconv.Atoi: parsing \"id\": invalid syntax"},
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
			if tt.args != nil && *tt.args != "" {
				path = fmt.Sprintf("/v1/issuer/%v", *tt.args)
				route = "/v1/issuer/{issuer_id}"
			} else {
				path = "/v1/issuer"
				route = "/v1/issuer"
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
			router.HandleFunc(route, handler.GetIssuer).Methods(method)
			router.ServeHTTP(writer, reader)

			// Asserts
			assert.NotNil(t, handler.service)
			if tt.expected != nil {
				expectedRespByte, _ := json.Marshal(tt.expected)
				var expectedResp issuerResponse
				_ = json.Unmarshal(expectedRespByte, &expectedResp)

				var resp issuerResponse
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
