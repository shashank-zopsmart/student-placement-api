package company

import (
	"encoding/json"
	"io"
	"net/http"
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
		w.WriteHeader(http.StatusMethodNotAllowed)
		response, _ := json.Marshal(entities.ResponseMessage{"Method not allowed"})
		w.Write(response)
	}
}

// Get handler to get company detail by ID
func (handler handler) Get(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(entities.ResponseMessage{"id param not present"})
		w.Write(response)
		return
	}

	result, err := handler.service.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response, _ := json.Marshal(entities.ResponseMessage{"Error: " + err.Error()})
		w.Write(response)
		return
	}
	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(result)
	w.Write(response)
}

// Create handler to create a new company
func (handler handler) Create(w http.ResponseWriter, req *http.Request) {
	reqBody, _ := io.ReadAll(req.Body)
	var company entities.Company
	json.Unmarshal(reqBody, &company)

	if company.Name == "" || company.Category == "" {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(entities.ResponseMessage{"Error: Name and Category required"})
		w.Write(response)
		return
	}

	result, err := handler.service.Create(company)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(entities.ResponseMessage{"Error: " + err.Error()})
		w.Write(response)
		return
	}

	w.WriteHeader(http.StatusCreated)
	response, _ := json.Marshal(result)
	w.Write(response)
}

// Update handler to update a particular company
func (handler handler) Update(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(entities.ResponseMessage{"Error: ID required"})
		w.Write(response)
		return
	}

	reqBody, _ := io.ReadAll(req.Body)
	var company entities.Company
	json.Unmarshal(reqBody, &company)

	if company.Name == "" || company.Category == "" {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(entities.ResponseMessage{"Error: Name and Category required"})
		w.Write(response)
		return
	}

	company.ID = id
	result, err := handler.service.Update(company)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(entities.ResponseMessage{"Error: " + err.Error()})
		w.Write(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(result)
	w.Write(response)
}

// Delete handler to delete a particular company
func (handler handler) Delete(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(entities.ResponseMessage{"Error: ID required"})
		w.Write(response)
		return
	}

	err := handler.service.Delete(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response, _ := json.Marshal(entities.ResponseMessage{"Error: Company not found"})
		w.Write(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(entities.ResponseMessage{"Company Deleted"})
	w.Write(response)
}
