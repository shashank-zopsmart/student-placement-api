package company

import (
	"encoding/json"
	"student-placement-api/errors"

	"io"
	"net/http"
	"strings"
	"student-placement-api/entities"
	"student-placement-api/service"
)

type handler struct {
	service service.Company
}

// New factory function to return handler object and do dependency injection
func New(service service.Company) handler {
	return handler{service: service}
}

// Handler main handler for the /company endpoint
func (handler handler) Handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case http.MethodGet:
		handler.Get(w, req)
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

// Create handler to create a new company
func (handler handler) Create(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := errors.InvalidParams{"Invalid request body"}
		w.Write([]byte(err.Error()))
		return
	}

	var company entities.Company

	err = json.Unmarshal(reqBody, &company)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := errors.InvalidParams{"Invalid request body"}
		w.Write([]byte(err.Error()))
		return
	}

	if company.Name == "" || company.Category == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := errors.MissisngParam{[]string{"ID", "Category"}}
		w.Write([]byte(err.Error()))
		return
	}

	result, err := handler.service.Create(ctx, company)
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

// Get handler to get company detail by ID
func (handler handler) Get(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id := req.URL.Query().Get("id")
	id = strings.TrimSpace(id)
	if id == "" {
		err := errors.MissisngParam{Params: []string{"id"}}
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	result, err := handler.service.GetByID(ctx, id)
	if err != nil {
		err := errors.EntityNotFound{Entity: "Company"}
		w.WriteHeader(http.StatusBadRequest)
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

// Update handler to update a particular company
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

	var company entities.Company
	err = json.Unmarshal(reqBody, &company)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		err := errors.InvalidParams{"Invalid request body"}
		w.Write([]byte(err.Error()))
		return
	}

	if company.Name == "" || company.Category == "" {
		w.WriteHeader(http.StatusBadRequest)
		err := errors.MissisngParam{[]string{"ID", "Category"}}
		w.Write([]byte(err.Error()))
		return
	}

	company.ID = id
	result, err := handler.service.Update(ctx, company)
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

// Delete handler to delete a particular company
func (handler handler) Delete(w http.ResponseWriter, req *http.Request) {
	ctx := req.Context()

	id := req.URL.Query().Get("id")
	id = strings.TrimSpace(id)

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

	response, err := json.Marshal(entities.ResponseMessage{"Company Deleted"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
