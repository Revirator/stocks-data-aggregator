package main

import (
	"fmt"
	"net/http"

	"github.com/revirator/cfd/views"
)

type ServerError struct {
	StatusCode int
	Message    string
}

func NewServerError(ticker string, err error) *ServerError {
	if err.Error() == "stock missing" {
		return &ServerError{
			StatusCode: http.StatusNotFound,
			Message:    fmt.Sprintf("Company with ticker '%s' does not exist or is not listed on any of the US exchanges.", ticker),
		}
	} else {
		return &ServerError{
			StatusCode: http.StatusInternalServerError,
			Message:    "Something went wrong! Please try again later.",
		}
	}
}

func CustomErrorHandler(statusCode int, message string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		views.Error(message).Render(r.Context(), w)
	}
}
