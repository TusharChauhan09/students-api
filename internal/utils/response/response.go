package response

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

const (
	StatusOk = "Ok"  
	StatusError = "Error"
)

func WriteJson(res http.ResponseWriter,status int ,data any) error {
	
	res.Header().Set("Content-Type","application/json")
	res.WriteHeader(status)

	return json.NewEncoder(res).Encode(data)
}


type Response struct {
	Status string `json:"status"`  // change into small cases in json
	Error string `json:"error"` 
}

func GeneralError (err error) Response {
	return Response {
		Status: StatusError,
		Error: err.Error(),
	}	
}


func ValidationError(errs validator.ValidationErrors) Response {
	var errMsgs []string

	for _,err := range errs {
		switch err.ActualTag() {
		case "required" :   // we have defined required in the struct that we are validating
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is required field",err.Field()))
		default:
			errMsgs = append(errMsgs, fmt.Sprintf("field %s is invalid",err.Field()))
		}
	}

	return Response{
		Status: StatusError,
		Error: strings.Join(errMsgs,", "),
	}
}