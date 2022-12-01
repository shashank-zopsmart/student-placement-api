package student

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"student-placement-api/entities"
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

		if id == "" && name == "" && branch == "" && includeCompany == "" {
			w.WriteHeader(http.StatusBadRequest)
			response, err := json.Marshal(entities.ErrorResponseMessage{"Error: Either ID or name, branch and " +
				"includeCompany required"})
			if err != nil {
				w.Write([]byte(fmt.Sprintf("Error: %v", err)))
				return
			}

			w.Write(response)
			return
		}

		if id != "" {
			handler.GetByID(w, req)
		} else {
			if name == "" || branch == "" || includeCompany == "" {
				w.WriteHeader(http.StatusBadRequest)
				response, err := json.Marshal(entities.ErrorResponseMessage{"Error: name, company and " +
					"includeCompany required"})
				if err != nil {
					w.Write([]byte(fmt.Sprintf("Error: %v", err)))
					return
				}

				w.Write(response)
			} else {
				handler.Get(w, req)
			}
		}
	case http.MethodPost:
		handler.Create(w, req)
	case http.MethodPut:
		handler.Update(w, req)
	case http.MethodDelete:
		handler.Delete(w, req)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		response, err := json.Marshal(entities.ErrorResponseMessage{"Method not allowed"})
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}

		w.Write(response)
	}
}

// Create handler to create new student
func (handler handler) Create(w http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response, err := json.Marshal(entities.ErrorResponseMessage{"Error: " + err.Error()})
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}

		w.Write(response)
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
		response, err := json.Marshal(entities.ErrorResponseMessage{"Error: Name, DOB, Branch, Phone, Company.ID, " +
			"Status required"})
		if err != nil {
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}
		w.Write(response)
		return
	}

	result, err := handler.service.Create(student)
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

	w.WriteHeader(http.StatusCreated)
	w.Write(response)
}

// Get handler to get all students
func (handler handler) Get(w http.ResponseWriter, req *http.Request) {
	name := strings.TrimSpace(req.URL.Query().Get("name"))
	branch := strings.TrimSpace(req.URL.Query().Get("branch"))
	includeCompany := strings.TrimSpace(req.URL.Query().Get("includeCompany"))

	if includeCompany == "" {
		includeCompany = "false"
	}

	includeCompanyFlag, err := strconv.ParseBool(includeCompany)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response, err := json.Marshal(entities.ErrorResponseMessage{"Error: " + err.Error()})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}

		w.Write(response)
		return
	}

	result, err := handler.service.Get(name, branch, includeCompanyFlag)
	if err != nil {
		response, err := json.Marshal(entities.ErrorResponseMessage{"Student not found"})
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

// GetByID handler ot get student by ID
func (handler handler) GetByID(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimSpace(req.URL.Query().Get("id"))

	result, err := handler.service.GetByID(id)
	if err != nil {
		response, err := json.Marshal(entities.ErrorResponseMessage{"Student not found"})
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

// Update handler to update a particular student
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

	var student entities.Student

	err = json.Unmarshal(reqBody, &student)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	if student.Name == "" || student.DOB == "" || student.Branch == "" || student.Phone == "" ||
		student.Company.ID == "" || student.Status == "" {

		response, err := json.Marshal(entities.ErrorResponseMessage{"Error: Name, DOB, Branch, Phone, Company.ID, " +
			"Status required"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}

		w.WriteHeader(http.StatusBadRequest)
		w.Write(response)
		return
	}

	student.ID = id
	result, err := handler.service.Update(student)
	if err != nil {
		if err == sql.ErrNoRows {
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

// Delete handler to delete a particular student
func (handler handler) Delete(w http.ResponseWriter, req *http.Request) {
	id := strings.TrimSpace(req.URL.Query().Get("id"))

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

		response, err := json.Marshal(entities.ErrorResponseMessage{"Error: Student not found"})
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte(fmt.Sprintf("Error: %v", err)))
			return
		}

		w.WriteHeader(http.StatusNotFound)
		w.Write(response)
		return
	}

	response, err := json.Marshal(entities.ErrorResponseMessage{"Student Deleted"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(fmt.Sprintf("Error: %v", err)))
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(response)
}
