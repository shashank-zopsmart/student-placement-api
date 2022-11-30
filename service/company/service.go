package company

import (
	"errors"
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
func (service Service) GetByID(id string) (entities.Company, error) {
	return service.store.GetByID(id)
}

// Create service to create a new company
func (service Service) Create(company entities.Company) (entities.Company, error) {
	switch company.Category {
	case "MASS", "DREAM IT", "OPEN DREAM", "CORE":
		return service.store.Create(company)
	default:
		return entities.Company{}, errors.New("invalid category")
	}
}

// Update service to update a particular company
func (service Service) Update(company entities.Company) (entities.Company, error) {
	_, err := service.store.GetByID(company.ID)
	if err != nil {
		return entities.Company{}, err
	}
	switch company.Category {
	case "MASS", "DREAM IT", "OPEN DREAM", "CORE":
		return service.store.Update(company)
	default:
		return entities.Company{}, errors.New("invalid category")
	}
}

// Delete service to delete a particular company
func (service Service) Delete(id string) error {
	_, err := service.store.GetByID(id)
	if err != nil {
		return err
	}
	return service.store.Delete(id)
}
