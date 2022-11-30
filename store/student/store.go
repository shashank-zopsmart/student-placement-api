package student

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

// GetById store to get a student by ID
func (store store) GetById(id string) (entities.Student, error) {
	return entities.Student{}, nil
}

// Get store to get all student or search student by name and branch
func (store store) Get(name string, branch string, includeCompany bool) ([]entities.Student, error) {
	return []entities.Student{}, nil
}

// Create store to create a new student
func (store store) Create(student entities.Student) (entities.Student, error) {
	return entities.Student{}, nil
}

// Update store to update a particular student
func (store store) Update(student entities.Student) (entities.Student, error) {
	return entities.Student{}, nil
}

// Delete store to delete a particular student
func (store store) Delete(id string) (entities.Student, error) {
	return entities.Student{}, nil
}

// GetCompany store to get company's detail a particular company id
func (store store) GetCompany(id string) (entities.Company, error) {
	return entities.Company{}, nil
}
