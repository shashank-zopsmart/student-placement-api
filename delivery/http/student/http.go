package student

import (
	"net/http"
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

}

// Get handler to get all students
func (handler Handler) Get(w http.ResponseWriter, req *http.Request) {

}

// GetById handler ot get student by ID
func (handler Handler) GetById(w http.ResponseWriter, req *http.Request) {

}

// Create handler to create new student
func (handler Handler) Create(w http.ResponseWriter, req *http.Request) {

}

// Update handler to update a particular student
func (handler Handler) Update(w http.ResponseWriter, req *http.Request) {

}

// Delete handler to delete a particular student
func (handler Handler) Delete(w http.ResponseWriter, req *http.Request) {

}
