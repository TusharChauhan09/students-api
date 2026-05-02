package student

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/TusharChauhan09/students-api/internal/storage"
	"github.com/TusharChauhan09/students-api/internal/types"
	"github.com/TusharChauhan09/students-api/internal/utils/response"
	"github.com/go-playground/validator/v10"
)

// storage storage.Storage : for the interface working
func New(storage storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		var student types.Student

		err := json.NewDecoder(req.Body).Decode(&student); 
		if errors.Is(err,io.EOF){
			response.WriteJson(res, http.StatusBadRequest, fmt.Errorf("empty body"))
			return
		}

		if err != nil {
			response.WriteJson(res,http.StatusBadRequest, response.GeneralError(err))
			return
		}

		// ! request validation  // go-playground
		if err := validator.New().Struct(student); err != nil{
			validateErrs := err.(validator.ValidationErrors)  // typecaste 
			response.WriteJson(res,http.StatusBadRequest,response.ValidationError(validateErrs))
			return
		}


		lastId, err := storage.CreateStudent(
			student.Name,
			student.Email,
			student.Age,
		)

		if err != nil {
			response.WriteJson(res, http.StatusInternalServerError,err)
			return
		}

		slog.Info("student created", slog.String("userId",fmt.Sprint(lastId)))

		response.WriteJson(res,http.StatusCreated, map[string]int64{"id" : lastId})

	}
}


func GetById(storage storage.Storage) http.HandlerFunc {
	return func(res http.ResponseWriter, req *http.Request) {
		id := req.PathValue("id")
		slog.Info("getting student", slog.String("id",id))

		intId,err := strconv.ParseInt(id, 10, 64)
		if err != nil {
			response.WriteJson(res,http.StatusBadRequest,response.GeneralError(err))
			return
		}

		student , err := storage.GetStudentById(intId)
		if err != nil {
			response.WriteJson(res, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}
		
		response.WriteJson(res, http.StatusOK, student) 
	}
}


func GetList(storage storage.Storage) http.HandlerFunc {
	return func (res http.ResponseWriter, req *http.Request) {
		students , err := storage.GetStudents()
		if err != nil {
			response.WriteJson(res, http.StatusInternalServerError, response.GeneralError(err))
			return 
		}

		response.WriteJson(res, http.StatusOK, students) 
	}
}