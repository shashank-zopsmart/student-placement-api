package student

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

// GetByID store to get a student by ID
func (store Store) GetByID(id string) (entities.Student, error) {
	return entities.Student{}, nil
}

// Get store to get all student or search student by name and branch
func (store Store) Get(name string, branch string, includeCompany bool) ([]entities.Student, error) {
	return []entities.Student{}, nil
}

// Create store to create a new student
func (store Store) Create(student entities.Student) error {
	return nil
}

// Update store to update a particular student
func (store Store) Update(student entities.Student) error {
	return nil
}
