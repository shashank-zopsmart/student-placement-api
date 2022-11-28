package company

import (
	"student-placement-api/entities"
	"student-placement-api/store"
)

type Service struct {
	store store.Company
}

// New factory function to return service object and do dependency injection
func New(store store.Company) Service {
	return Service{store: store}
}

// GetByID service to get a company by ID
func (service Service) GetByID(id string) entities.Company {
	return entities.Company{}
}

// Create service to create a new company
func (service Service) Create(company entities.Company) error {
	return nil
}

// Update service to update a particular company
func (service Service) Update(company entities.Company) error {
	return nil
}

// Delete service to delete a particular company
func (service Service) Delete(id string) error {
	return nil
}
