package student

import (
	"encoding/json"
	"io"
	"net/http"
	"strconv"
	"student-placement-api/entities"
	"student-placement-api/service"
)

type Handler struct {
	service service.Student
}

// New factory function to return handler object and do dependency injection
func New(service service.Student) Handler {
	return Handler{service: service}
}

// Handler main handler for /student endpoint
func (handler Handler) Handler(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	switch req.Method {
	case http.MethodGet:
		id := req.URL.Query().Get("id")
		name := req.URL.Query().Get("name")
		branch := req.URL.Query().Get("branch")
		includeCompany := req.URL.Query().Get("includeCompany")

		if id == "" && name == "" && branch == "" && includeCompany == "" {
			w.WriteHeader(http.StatusBadRequest)
			response, _ := json.Marshal(entities.ResponseMessage{"Error: Either ID or name, company and " +
				"includeCompany required"})
			w.Write(response)
		}

		if id != "" {
			handler.GetByID(w, req)
		} else {
			if name == "" || branch == "" || includeCompany == "" {
				w.WriteHeader(http.StatusBadRequest)
				response, _ := json.Marshal(entities.ResponseMessage{"Error: name, company and " +
					"includeCompany required"})
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
		response, _ := json.Marshal(entities.ResponseMessage{"Method not allowed"})
		w.Write(response)
	}
}

// Get handler to get all students
func (handler Handler) Get(w http.ResponseWriter, req *http.Request) {
	name := req.URL.Query().Get("name")
	branch := req.URL.Query().Get("branch")
	includeCompany := req.URL.Query().Get("includeCompany")

	if includeCompany == "" {
		includeCompany = "false"
	}
	includeCompanyFlag, err := strconv.ParseBool(includeCompany)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response, _ := json.Marshal(entities.ResponseMessage{"Error: " + err.Error()})
		w.Write(response)
		return
	}

	result, err := handler.service.Get(name, branch, includeCompanyFlag)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response, _ := json.Marshal(result)
		w.Write(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(result)
	w.Write(response)
}

// GetByID handler ot get student by ID
func (handler Handler) GetByID(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")

	result, err := handler.service.GetByID(id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response, _ := json.Marshal(result)
		w.Write(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(result)
	w.Write(response)
}

// Create handler to create new student
func (handler Handler) Create(w http.ResponseWriter, req *http.Request) {
	reqBody, err := io.ReadAll(req.Body)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		response, _ := json.Marshal(entities.ResponseMessage{"Error: " + err.Error()})
		w.Write(response)
		return
	}

	var student entities.Student
	json.Unmarshal(reqBody, &student)

	if student.Name == "" || student.DOB == "" || student.Branch == "" || student.Phone == "" ||
		student.Company.ID == "" || student.Status == "" {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(entities.ResponseMessage{"Error: Name, DOB, Branch, Phone, Company.ID, " +
			"Status required"})
		w.Write(response)
		return
	}

	result, err := handler.service.Create(student)
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

// Update handler to update a particular student
func (handler Handler) Update(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(entities.ResponseMessage{"Error: ID required"})
		w.Write(response)
		return
	}

	reqBody, _ := io.ReadAll(req.Body)
	var student entities.Student
	json.Unmarshal(reqBody, &student)

	if student.Name == "" || student.DOB == "" || student.Branch == "" || student.Phone == "" ||
		student.Company.ID == "" || student.Status == "" {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(entities.ResponseMessage{"Error: Name, DOB, Branch, Phone, Company.ID, " +
			"Status required"})
		w.Write(response)
		return
	}

	student.ID = id
	result, err := handler.service.Update(student)
	if err != nil {
		if err.Error() == "student not found" {
			w.WriteHeader(http.StatusNotFound)
			response, _ := json.Marshal(entities.ResponseMessage{"Error: " + err.Error()})
			w.Write(response)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(entities.ResponseMessage{"Error: " + err.Error()})
		w.Write(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(result)
	w.Write(response)
}

// Delete handler to delete a particular student
func (handler Handler) Delete(w http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get("id")
	if id == "" {
		w.WriteHeader(http.StatusBadRequest)
		response, _ := json.Marshal(entities.ResponseMessage{"Error: ID required"})
		w.Write(response)
		return
	}

	_, err := handler.service.Delete(id)

	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		response, _ := json.Marshal(entities.ResponseMessage{"Error: Student not found"})
		w.Write(response)
		return
	}

	w.WriteHeader(http.StatusOK)
	response, _ := json.Marshal(entities.ResponseMessage{"Student Deleted"})
	w.Write(response)
}
