package company

import (
	"encoding/json"
	"fmt"
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
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte(fmt.Sprintf("Method not allowed")))
	}
}

// Create handler to create a new company
func (handler handler) Create(w http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}
	var company entities.Company

	err = json.Unmarshal(reqBody, &company)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	if company.Name == "" || company.Category == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: Name and Category required")))
		return
	}

	result, err := handler.service.Create(company)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	response, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

// Get handler to get company detail by ID
func (handler handler) Get(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	id = strings.TrimSpace(id)

	if id == "" {
		response, err := json.Marshal(entities.ErrorResponseMessage{"id param not present"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	result, err := handler.service.GetByID(id)
	if err != nil {
		response, err := json.Marshal(entities.ErrorResponseMessage{"Company not found"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}

	response, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Update handler to update a particular company
func (handler handler) Update(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	id = strings.TrimSpace(id)

	if id == "" {
		response, err := json.Marshal(entities.ErrorResponseMessage{"Error: ID required"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	var company entities.Company
	err = json.Unmarshal(reqBody, &company)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	if company.Name == "" || company.Category == "" {
		response, err := json.Marshal(entities.ErrorResponseMessage{"Error: Name and Category required"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	company.ID = id
	result, err := handler.service.Update(company)
	if err != nil {
		response, err := json.Marshal(entities.ErrorResponseMessage{"Error: " + err.Error()})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	response, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}

// Delete handler to delete a particular company
func (handler handler) Delete(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	id = strings.TrimSpace(id)

	if id == "" {
		response, err := json.Marshal(entities.ErrorResponseMessage{"Error: ID required"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	err := handler.service.Delete(id)
	if err != nil {
		response, err := json.Marshal(entities.ErrorResponseMessage{"Error: " + err.Error()})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}

	response, err := json.Marshal(entities.ErrorResponseMessage{"Company Deleted"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
