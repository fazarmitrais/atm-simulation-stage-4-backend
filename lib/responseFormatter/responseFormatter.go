package responseFormatter

import (
	"encoding/json"
	"net/http"
)

type ResponseFormatter struct {
	StatusCode int
	Message    string
	IsError    bool
}

func New(statusCode int, message string, isError bool) *ResponseFormatter {
	return &ResponseFormatter{StatusCode: statusCode, Message: message, IsError: isError}
}

func (r *ResponseFormatter) ReturnAsJson(w http.ResponseWriter) {
	w.Header().Add("content-type", "application/json")
	if r == nil {
		r = &ResponseFormatter{StatusCode: http.StatusOK, Message: "OK"}
	}
	w.WriteHeader(r.StatusCode)
	json.NewEncoder(w).Encode(r)
}
