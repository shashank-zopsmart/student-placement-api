package student

import (
	"student-placement-api/entities"
	"student-placement-api/store"
)

type Service struct {
	store store.Student
}

// New factory function to return service object and do dependency injection
func New(store store.Student) Service {
	return Service{store: store}
}

// Get service to get all student or search student by name and branch
func (service Service) Get(name string, branch string, includeCompany bool) ([]entities.Student, error) {
	return []entities.Student{}, nil
}

// GetByID service to get a student by ID
func (service Service) GetByID(id string) (entities.Student, error) {
	return entities.Student{}, nil
}

// Create to create a new student
func (service Service) Create(company entities.Student) error {
	return nil
}

// Update service to update a particular student
func (service Service) Update(company entities.Student) error {
	return nil
}

// Delete service to delete a particular student
func (service Service) Delete(id string) error {
	return nil
}
