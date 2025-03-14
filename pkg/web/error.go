package web

import (
	"encoding/json"
	"errors"
	errs "github.com/esvarez/go-api/pkg/error"
	"net/http"
)

var (
	InternalServerError   = AppError{StatusCode: http.StatusInternalServerError, Type: "internal_server_error", Message: "Internal server error"}
	ResourceNotFoundError = AppError{StatusCode: http.StatusNotFound, Type: "resource_not_found", Message: "Resource not found"}
	BadRequestError       = AppError{StatusCode: http.StatusBadRequest, Type: "bad_request", Message: "Bad request"}
)

type AppError struct {
	StatusCode int    `json:"status_code"`
	Type       string `json:"type"`
	Message    string `json:"message"`
}

func (e AppError) Send(w http.ResponseWriter) error {
	statusCode := e.StatusCode

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(e)
}

func ErrorResponse(err error) AppError {
	if err == nil {
		return InternalServerError
	}
	switch {
	case errors.Is(err, errs.ErrNotFound):
		return ResourceNotFoundError
	default:
		return InternalServerError
	}
}
