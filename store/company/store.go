package company

import (
	"context"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/uuid"
	"student-placement-api/entities"
)

type store struct {
	db *sql.DB
}

// New factory function to return store object and do dependency injection
func New(db *sql.DB) store {
	return store{db}
}

// Create store to create a new company
func (store store) Create(ctx context.Context, company entities.Company) (entities.Company, error) {
	company.ID = uuid.New().String()
	query := "INSERT INTO companies (id, name, category) VALUES(?, ?, ?)"

	_, err := store.db.Exec(query, company.ID, company.Name, company.Category)
	if err != nil {
		return entities.Company{}, err
	}

	return company, nil
}

// GetByID store to get a company by ID
func (store store) GetByID(ctx context.Context, id string) (entities.Company, error) {
	query := "SELECT * FROM companies WHERE id=?"

	var company entities.Company
	row := store.db.QueryRow(query, id)
	err := row.Scan(&company.ID, &company.Name, &company.Category)
	if err != nil {
		return entities.Company{}, err
	}

	return company, nil
}

// Update store to update a particular company
func (store store) Update(ctx context.Context, company entities.Company) (entities.Company, error) {
	query := "UPDATE companies SET name=?, category=? WHERE id=?"

	_, err := store.db.Exec(query, company.Name, company.Category, company.ID)
	if err != nil {
		return entities.Company{}, err
	}

	return company, nil
}

// Delete store to delete a particular company
func (store store) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM companies WHERE id=?"
	_, err := store.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}
