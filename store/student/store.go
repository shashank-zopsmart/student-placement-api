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

// GetById store to get a student by ID
func (store Store) GetById(id string) (entities.Student, error) {
	return entities.Student{}, nil
}

// Get store to get all student or search student by name and branch
func (store Store) Get(name string, branch string, includeCompany bool) ([]entities.Student, error) {
	return []entities.Student{}, nil
}

// Create store to create a new student
func (store Store) Create(student entities.Student) (entities.Student, error) {
	return entities.Student{}, nil
}

// Update store to update a particular student
func (store Store) Update(student entities.Student) (entities.Student, error) {
	return entities.Student{}, nil
}

// Delete store to delete a particular student
func (store Store) Delete(id string) (entities.Student, error) {
	return entities.Student{}, nil
}

// GetCompany store to get company's detail a particular company id
func (store Store) GetCompany(id string) (entities.Company, error) {
	return entities.Company{}, nil
}
