package student

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"student-placement-api/entities"
	"student-placement-api/errors"
	"student-placement-api/service"
)

type handler struct {
	service service.Student
}

// New factory function to return handler object and do dependency injection
func New(service service.Student) handler {
	return handler{service: service}
}

// Handler main handler for /student endpoint
func (handler handler) Handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case http.MethodGet:
		id := req.URL.Query().Get("id")
		name := req.URL.Query().Get("name")
		branch := req.URL.Query().Get("branch")
		includeCompany := req.URL.Query().Get("includeCompany")

		if id == "" && (name == "" || branch == "" || includeCompany == "") {
			err := errors.MissisngParam{Params: []string{"id", "name", "branch", "includeCompany"}}
			w.WriteHeader(http.StatusBadRequest)
			w.Write([]byte(err.Error()))
			return
		}

		if id != "" {
			handler.GetByID(w, req)
		} else {
			handler.Get(w, req)
		}
	case http.MethodPost:
		handler.Create(w, req)
	case http.MethodPut:
		handler.Update(w, req)
	case http.MethodDelete:
		handler.Delete(w, req)
	default:
		err := errors.HttpErrors{"Method not allowed"}
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(err.Error()))
	}
}

// Create handler to create new student
func (handler handler) Create(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := errors.InvalidParams{"Invalid request body"}
		w.Write([]byte(err.Error()))
		return
	}

	var student entities.Student
	err = json.Unmarshal(reqBody, &student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	if student.Name == "" || student.DOB == "" || student.Branch == "" || student.Phone == "" ||
		student.Company.ID == "" || student.Status == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := errors.MissisngParam{[]string{"name", "dob", "branch", "phone", "company_id", "status"}}
		w.Write([]byte(err.Error()))
		return
	}

	result, err := handler.service.Create(ctx, student)
	if err != nil {
		switch err.(type) {
		case errors.ConnDone, errors.InvalidParams:
			w.WriteHeader(http.StatusBadRequest)
		}
		w.Write([]byte(err.Error()))
		return
	}

	response, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

// Get handler to get all students
func (handler handler) Get(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	name := strings.TrimSpace(req.URL.Query().Get("name"))
	branch := strings.TrimSpace(req.URL.Query().Get("branch"))
	includeCompany := strings.TrimSpace(req.URL.Query().Get("includeCompany"))

	includeCompanyFlag, err := strconv.ParseBool(includeCompany)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err = errors.InvalidParams{"includeCompany should be either true or false"}
		w.Write([]byte(err.Error()))
		return
	}

	result, err := handler.service.Get(ctx, name, branch, includeCompanyFlag)
	if err != nil {
		err := errors.EntityNotFound{Entity: "Student"}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	response, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// GetByID handler ot get student by ID
func (handler handler) GetByID(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id := strings.TrimSpace(req.URL.Query().Get("id"))

	result, err := handler.service.GetByID(ctx, id)
	if err != nil {
		err := errors.EntityNotFound{Entity: "Student"}
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte(err.Error()))
		return
	}

	response, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Update handler to update a particular student
func (handler handler) Update(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id := req.URL.Query().Get("id")
	id = strings.TrimSpace(id)

	if id == "" {
		err := errors.MissisngParam{Params: []string{"id"}}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := errors.InvalidParams{"Invalid request body"}
		w.Write([]byte(err.Error()))
		return
	}

	var student entities.Student

	err = json.Unmarshal(reqBody, &student)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := errors.InvalidParams{"Invalid request body"}
		w.Write([]byte(err.Error()))
		return
	}

	if student.Name == "" || student.DOB == "" || student.Branch == "" || student.Phone == "" ||
		student.Company.ID == "" || student.Status == "" {

		w.WriteHeader(http.StatusBadRequest)
		err := errors.MissisngParam{[]string{"name", "dob", "branch", "phone", "company_id", "status"}}
		w.Write([]byte(err.Error()))
		return
	}

	student.ID = id
	result, err := handler.service.Update(ctx, student)
	if err != nil {
		switch err.(type) {
		case errors.ConnDone, errors.InvalidParams:
			w.WriteHeader(http.StatusBadRequest)
		case errors.EntityNotFound:
			w.WriteHeader(http.StatusNotFound)
		}

		w.Write([]byte(err.Error()))
		return
	}

	response, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Delete handler to delete a particular student
func (handler handler) Delete(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id := strings.TrimSpace(req.URL.Query().Get("id"))

	if id == "" {
		err := errors.MissisngParam{Params: []string{"id"}}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	err := handler.service.Delete(ctx, id)

	if err != nil {
		switch err.(type) {
		case errors.ConnDone, errors.InvalidParams:
			w.WriteHeader(http.StatusBadRequest)
		case errors.EntityNotFound:
			w.WriteHeader(http.StatusNotFound)
		}
		w.Write([]byte(err.Error()))
		return
	}

	response, err := json.Marshal(entities.ResponseMessage{"Student Deleted"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusNoContent)
	w.Write(response)
}
