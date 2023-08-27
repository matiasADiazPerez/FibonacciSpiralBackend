package utils

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"net/http"
	"spiralmatrix/config"

	"github.com/go-playground/validator/v10"
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

func GetBody[M any](r *http.Request, model M) (M, ErrorWrapper) {

	err := json.NewDecoder(r.Body).Decode(&model)
	if err != nil {
		return model, NewErrorWrapper(config.BODY, 0, err)
	}
	validate := validator.New()
	if err := validate.Struct(model); err != nil {
		return model, NewErrorWrapper(config.BODY, http.StatusBadRequest, err)
	}
	return model, ErrorWrapper{}

}

func CreateResponse(message string, payload any, w http.ResponseWriter) {
	response := Response{
		Message: message,
		Payload: payload,
	}
	jData, err := json.Marshal(response)
	if err != nil {
		HandleError(NewErrorWrapper("Error creating the response", 0, err), w)
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

func HashPassword(password string) (string, error) {
	hasher := sha256.New()
	_, err := hasher.Write([]byte(password))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hasher.Sum(nil)), nil

}
