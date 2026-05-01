package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"

	"github.com/TusharChauhan09/students-api/internal/storage"
	"github.com/TusharChauhan09/students-api/internal/types"
	"github.com/TusharChauhan09/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

func New(storage storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var student types.Student

		err := json.NewDecoder(req.Body).Decode(&student); 
		if errors.Is(err,io.EOF){
			response.WriteJson(res, http.StatusBadRequest, fmt.Errorf("empty body"))
			return
		}

		if err != nil {
			response.WriteJson(res,http.StatusBadRequest, response.GenralError(err))
		}

		// ! request validation  // go-playground
		if err := validator.New().Struct(student); err != nil{
			validateErrs := err.(validator.ValidationErrors)  // typecaste 
			response.WriteJson(res,http.StatusBadRequest,response.ValidationError(validateErrs))
			return
		}


		slog.Info("student created")

		response.WriteJson(res,http.StatusCreated, map[string]string{"success" : "ok"})

	}
}