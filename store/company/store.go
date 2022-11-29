package company

import (
	"database/sql"
	"student-placement-api/entities"
)

type Store struct {
	db *sql.DB
}

// New factory function to return store object and do dependency injection
func New(db *sql.DB) Store {
	return Store{db}
}

// GetByID store to get a company by ID
func (store Store) GetByID(id string) (entities.Company, error) {
	return entities.Company{}, nil
}

// Create store to create a new company
func (store Store) Create(company entities.Company) (entities.Student, error) {
	return entities.Student{}, nil
}

// Update store to update a particular company
func (store Store) Update(company entities.Company) (entities.Student, error) {
	return entities.Student{}, nil
}

// Delete store to delete a particular company
func (store Store) Delete(id string) (entities.Student, error) {
	return entities.Student{}, nil
}
