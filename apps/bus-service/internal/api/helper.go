package api

import (
	"encoding/json"
	"errors"
	"github.com/perkzen/mbus/apps/bus-service/internal/errs"
	"github.com/perkzen/mbus/apps/bus-service/internal/utils"
	"log/slog"
	"net/http"
	"strconv"
)

type HandlerFunc func(w http.ResponseWriter, r *http.Request) error

func MakeHandlerFunc(fn HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := fn(w, r); err != nil {
			var apiErr errs.APIError
			ok := errors.As(err, &apiErr)
			if !ok {
				slog.Error(err.Error())
				apiErr = errs.InternalServerError()
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
		return "", errs.BadRequestError("missing query parameter: " + key)
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
