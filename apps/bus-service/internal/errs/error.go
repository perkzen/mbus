package errs

import (
	"fmt"
	"net/http"
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

func BusStationNotFoundError(id int) APIError {
	return NotFoundError(fmt.Sprintf("Bus station with ID %d does not exist", id))
}
