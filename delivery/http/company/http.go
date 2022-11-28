package company

import (
	"net/http"
	"student-placement-api/service"
)

type Handler struct {
	service service.Company
}

// New factory function to return handler object and do dependency injection
func New(service service.Company) Handler {
	return Handler{service: service}
}

// Handler main handler for the /company endpoint
func (handler Handler) Handler(w http.ResponseWriter, req *http.Request) {

}

// Get handler to get company detail by ID
func (handler Handler) Get(w http.ResponseWriter, req *http.Request) {

}

// Create handler to create a new company
func (handler Handler) Create(w http.ResponseWriter, req *http.Request) {

}

// Update handler to update a particular company
func (handler Handler) Update(w http.ResponseWriter, req *http.Request) {

}

// Delete handler to delete a particular company
func (handler Handler) Delete(w http.ResponseWriter, req *http.Request) {

}
