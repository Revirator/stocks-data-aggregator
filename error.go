package main

import "net/http"

type ServerError struct {
	StatusCode int
	Message    string
}

func InternalServerError() *ServerError {
	return &ServerError{
		StatusCode: http.StatusInternalServerError,
		Message:    "Something went wrong! Please try again later.",
	}
}

// TODO: update template
func CustomerErrorHandler(statusCode int, message string) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		WriteHTML(writer, statusCode, "error.html", message)
	}
}
