package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Message string
	Payload any
}

type ErrorWrapper struct {
	Message string
	Error   error
	Code    int
}

func CreateResponse(message string, payload any, w http.ResponseWriter) {
	response := Response{
		Message: message,
		Payload: payload,
	}
	jData, err := json.Marshal(response)
	if err != nil {
		HandleError(NewErrorWrapper("Error creating the response", http.StatusInternalServerError, err), w)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jData)
}

func NewErrorWrapper(message string, code int, err error) ErrorWrapper {
	if code == 0 {
		code = http.StatusInternalServerError
	}
	return ErrorWrapper{
		Message: message,
		Error:   err,
		Code:    code,
	}
}

func HandleError(err ErrorWrapper, w http.ResponseWriter) {
	response := Response{
		Message: err.Message,
		Payload: err.Error.Error(),
	}
	jData, jsonErr := json.Marshal(response)
	if jsonErr != nil {
		jData = []byte(fmt.Sprintf("Unexpected Err: %s", jsonErr.Error()))
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(err.Code)
	w.Write(jData)

}
