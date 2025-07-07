package api

import (
	"encoding/json"
	"errors"
	"github.com/perkzen/mbus/bus-service/internal/utils"
	"net/http"
	"strconv"
)

type APIError struct {
	StatusCode int    `json:"statusCode"`
	Message    string `json:"message"`
}

func NewAPIError(statusCode int, message string) APIError {
	return APIError{
		StatusCode: statusCode,
		Message:    message,
	}
}

func (e APIError) Error() string {
	return e.Message
}

func InternalServerError() APIError {
	return NewAPIError(http.StatusInternalServerError, "Internal Server Error")
}

func BadRequestError(message string) APIError {
	return NewAPIError(http.StatusBadRequest, message)
}

func NotFoundError(message string) APIError {
	return NewAPIError(http.StatusNotFound, message)
}

type APIFunc func(w http.ResponseWriter, r *http.Request) error

func MakeHandlerFunc(fn APIFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			var apiErr APIError
			ok := errors.As(err, &apiErr)
			if !ok {
				apiErr = InternalServerError()
			}
			_ = WriteJSON(w, apiErr.StatusCode, apiErr)

		}
	}
}

func WriteJSON(w http.ResponseWriter, status int, data any) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(data)
}

func QueryInt(r *http.Request, key string, defaultValue int) int {
	val := r.URL.Query().Get(key)
	if parsed, err := strconv.Atoi(val); err == nil {
		return parsed
	}
	return defaultValue
}

func QueryStr(r *http.Request, key string) (string, error) {
	val := r.URL.Query().Get(key)
	if val == "" {
		return "", BadRequestError("missing query parameter: " + key)
	}

	return val, nil

}

func QueryDateStr(r *http.Request, key string, defaultValue string) string {
	date, _ := QueryStr(r, key)

	if date != "" && !utils.ValidateDate(date) {
		return defaultValue
	}

	if date == "" {
		return defaultValue
	}

	return date
}
