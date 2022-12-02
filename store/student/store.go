package student

import (
	"context"
	"database/sql"
	"github.com/google/uuid"
	"student-placement-api/entities"
	"student-placement-api/errors"
)

type store struct {
	db *sql.DB
}

// New factory function to return store object and do dependency injection
func New(db *sql.DB) store {
	return store{db}
}

// Create store to create a new student
func (store store) Create(ctx context.Context, student entities.Student) (entities.Student, error) {
	student.ID = uuid.NewString()
	query := "INSERT INTO students (id, name, dob, phone, branch, company_id, status) VALUES(?, ?, ?, ?, ?, ?, ?)"
	_, err := store.db.Exec(query, student.ID, student.Name, student.DOB, student.Phone, student.Branch,
		student.Company.ID, student.Status)
	if err != nil {
		return entities.Student{}, errors.ConnDone{}
	}
	return student, nil
}

// GetById store to get a student by ID
func (store store) GetById(ctx context.Context, id string) (entities.Student, error) {
	query := "SELECT students.id AS id, students.name AS name, students.dob AS dob, students.phone AS phone, " +
		"students.branch AS branch, students.status AS status FROM students WHERE students.id=?"
	row := store.db.QueryRow(query, id)
	var student entities.Student

	if err := row.Scan(&student.ID, &student.Name, &student.DOB, &student.Phone, &student.Branch,
		&student.Status); err != nil {
		return entities.Student{}, errors.EntityNotFound{Entity: "Student"}
	}

	return student, nil
}

// Get store to get all student or search student by name and branch
func (store store) Get(ctx context.Context, name string, branch string, includeCompany bool) ([]entities.Student, error) {
	var students = make([]entities.Student, 0)

	if includeCompany == true {
		query := "SELECT students.id AS id, students.name AS name, students.dob AS dob, students.phone AS phone, " +
			"students.branch AS branch, companies.id AS companyID, companies.name AS companyName, " +
			"companies.category AS companyCategory, students.status AS status FROM students JOIN companies ON " +
			"students.company_id=companies.id WHERE students.name=? AND students.branch=?"

		rows, err := store.db.Query(query, name, branch)
		if err != nil {
			return []entities.Student{}, errors.ConnDone{}
		}

		for rows.Next() {
			var student entities.Student
			rows.Scan(&student.ID, &student.Name, &student.DOB, &student.Phone, &student.Branch, &student.Company.ID,
				&student.Company.Name, &student.Company.Category, &student.Status)
			students = append(students, student)
		}
	} else {
		query := "SELECT id, name, dob, phone, branch, status FROM students WHERE students.name=? AND students.branch=?"

		rows, err := store.db.Query(query, name, branch)
		if err != nil {
			return []entities.Student{}, errors.ConnDone{}
		}

		for rows.Next() {
			var student entities.Student
			rows.Scan(&student.ID, &student.Name, &student.DOB, &student.Phone, &student.Branch,
				&student.Status)
			students = append(students, student)
		}
	}

	if len(students) == 0 {
		return students, errors.EntityNotFound{Entity: "Student"}
	}

	return students, nil
}

// Update store to update a particular student
func (store store) Update(ctx context.Context, student entities.Student) (entities.Student, error) {
	query := "UPDATE students SET name=?, phone=?, dob=?, branch=?, company_id=?, status=? WHERE id=?"
	_, err := store.db.Exec(query, student.Name, student.Phone, student.DOB, student.Branch, student.Company.ID,
		student.Status, student.ID)

	if err != nil {
		return entities.Student{}, errors.ConnDone{}
	}
	return student, nil
}

// Delete store to delete a particular student
func (store store) Delete(ctx context.Context, id string) error {
	query := "DELETE FROM students WHERE id=?"
	_, err := store.db.Exec(query, id)
	if err != nil {
		return errors.ConnDone{}
	}
	return nil
}

// GetCompany store to get company's detail a particular company id
func (store store) GetCompany(ctx context.Context, id string) (entities.Company, error) {
	query := "SELECT * FROM companies WHERE id=?"

	var company entities.Company
	row := store.db.QueryRow(query, id)
	err := row.Scan(&company.ID, &company.Name, &company.Category)

	if err != nil {
		return entities.Company{}, errors.EntityNotFound{Entity: "Company"}
	}

	return company, nil
}
