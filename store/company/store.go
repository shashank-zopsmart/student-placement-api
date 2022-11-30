package company

import (
	"database/sql"
	"student-placement-api/entities"
)

type store struct {
	db *sql.DB
}

// New factory function to return store object and do dependency injection
func New(db *sql.DB) store {
	return store{db}
}

// GetByID store to get a company by ID
func (store store) GetByID(id string) (entities.Company, error) {
	return entities.Company{}, nil
}

// Create store to create a new company
func (store store) Create(company entities.Company) (entities.Student, error) {
	return entities.Student{}, nil
}

// Update store to update a particular company
func (store store) Update(company entities.Company) (entities.Student, error) {
	return entities.Student{}, nil
}

// Delete store to delete a particular company
func (store store) Delete(id string) error {
	return nil
}
