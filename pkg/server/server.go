package server

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
)

type errorResponse struct {
	Code     string
	Messages []string
}

// OK renders obj with 200 status code
func OK(w http.ResponseWriter, r *http.Request, obj interface{}) {
	Render(w, r, obj, http.StatusOK)
}

// OKNoContent empty body with 204 status code
func OKNoContent(w http.ResponseWriter) {
	w.WriteHeader(http.StatusNoContent)
}

// BadRequest response as json
func BadRequest(w http.ResponseWriter, r *http.Request, code string, messages ...string) {
	err := &errorResponse{
		Code:     code,
		Messages: messages,
	}
	Render(w, r, err, http.StatusBadRequest)
}

// InternalServerError response as json
func InternalServerError(w http.ResponseWriter, r *http.Request, err error) {
	Render(w, r, &errorResponse{
		Code:     "INTERNAL_SERVER_ERROR",
		Messages: []string{err.Error()},
	}, http.StatusInternalServerError)
}

// Render response as json
func Render(w http.ResponseWriter, r *http.Request, obj interface{}, status int) {
	js, err := json.Marshal(obj)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	w.Write(js)
}

// GetStringFromPath gets an string value from a named path part
func GetStringFromPath(r *http.Request, key, defaultValue string) string {
	str := mux.Vars(r)[key]

	if len(str) < 1 {
		return defaultValue
	}

	return str
}

// GetStringParam returns a string param from query string
func GetStringParam(r *http.Request, key, defaultValue string) string {
	keys, ok := r.URL.Query()[key]

	if !ok || len(keys[0]) < 1 {
		return defaultValue
	}

	return keys[0]
}
