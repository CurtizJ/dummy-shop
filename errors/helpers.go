package errors

import (
	"encoding/json"
	"net/http"
)

type jsonError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
}

func ReportAsJSON(w http.ResponseWriter, message string, status int) {
	jsonMessage, _ := json.Marshal(jsonError{UNKNOWN_ERROR, message})
	http.Error(w, string(jsonMessage), status)
}

func ReportErrorAsJSON(w http.ResponseWriter, err error) {
	if appError, ok := err.(ApplicationError); ok {
		jsonMessage, _ := json.Marshal(jsonError{appError.Code(), err.Error()})
		http.Error(w, string(jsonMessage), appError.GetHTTPStatus())

	} else {
		ReportAsJSON(w, err.Error(), http.StatusInternalServerError)
	}
}
