package main

import (
	"fmt"
	"net/http"

	"github.com/revirator/cfd/view"
)

type ServerError struct {
	StatusCode int
	Message    string
}

const MISSING_MESSAGE = "Company with ticker '%s' does not exist or is not listed on any of the US exchanges."

func NewServerError(ticker string, err error) *ServerError {
	if err.Error() == "company missing" {
		return &ServerError{
			StatusCode: http.StatusNotFound,
			Message:    fmt.Sprintf(MISSING_MESSAGE, ticker),
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
		view.Error(message).Render(r.Context(), w)
	}
}
