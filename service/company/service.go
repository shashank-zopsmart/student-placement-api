package company

import (
	"context"
	"errors"
	"student-placement-api/entities"
	"student-placement-api/store"
)

type service struct {
	store store.Company
}

// New factory function to return service object and do dependency injection
func New(store store.Company) service {
	return service{store: store}
}

// Create service to create a new company
func (service service) Create(ctx context.Context, company entities.Company) (entities.Company, error) {
	switch company.Category {
	case "MASS", "DREAM IT", "OPEN DREAM", "CORE":
		return service.store.Create(ctx, company)
	default:
		return entities.Company{}, errors.New("invalid category")
	}
}

// GetByID service to get a company by ID
func (service service) GetByID(ctx context.Context, id string) (entities.Company, error) {
	return service.store.GetByID(ctx, id)
}

// Update service to update a particular company
func (service service) Update(ctx context.Context, company entities.Company) (entities.Company, error) {
	_, err := service.store.GetByID(ctx, company.ID)
	if err != nil {
		return entities.Company{}, err
	}
	switch company.Category {
	case "MASS", "DREAM IT", "OPEN DREAM", "CORE":
		return service.store.Update(ctx, company)
	default:
		return entities.Company{}, errors.New("invalid category")
	}
}

// Delete service to delete a particular company
func (service service) Delete(ctx context.Context, id string) error {
	_, err := service.store.GetByID(ctx, id)
	if err != nil {
		return err
	}
	return service.store.Delete(ctx, id)
}
