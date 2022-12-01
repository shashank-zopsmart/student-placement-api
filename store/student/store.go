package student

import (
	"database/sql"
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

// GetById store to get a student by ID
func (store store) GetById(id string) (entities.Student, error) {
	if err := store.db.Ping(); err != nil {
		return entities.Student{}, sql.ErrConnDone
	}

	query := "SELECT student.id AS id, student.name AS name, student.dob AS dob, student.phone AS phone, " +
		"student.branch AS branch, student.status AS status FROM student WHERE student.id=?"
	row := store.db.QueryRow(query, id)
	var student entities.Student

	if err := row.Scan(&student.ID, &student.Name, &student.DOB, &student.Phone, &student.Branch,
		&student.Status); err != nil {
		return entities.Student{}, err
	}

	return student, nil
}

// Get store to get all student or search student by name and branch
func (store store) Get(name string, branch string, includeCompany bool) ([]entities.Student, error) {
	if err := store.db.Ping(); err != nil {
		return []entities.Student{}, sql.ErrConnDone
	}

	var students = make([]entities.Student, 0)

	if includeCompany == true {
		query := "SELECT student.id AS id, student.name AS name, student.dob AS dob, student.phone AS phone, " +
			"student.branch AS branch, company.id AS companyID, company.name AS companyName, " +
			"company.category AS companyCategory, student.status AS status FROM student JOIN company ON " +
			"student.companyID=company.id WHERE student.name=? AND student.branch=?"

		rows, err := store.db.Query(query, name, branch)
		if err != nil {
			return []entities.Student{}, err
		}

		for rows.Next() {
			var student entities.Student
			rows.Scan(&student.ID, &student.Name, &student.DOB, &student.Phone, &student.Branch, &student.Company.ID,
				&student.Company.Name, &student.Company.Category, &student.Status)
			students = append(students, student)
		}
	} else {
		query := "SELECT * FROM student WHERE student.name=? AND student.branch=?"

		rows, err := store.db.Query(query, name, branch)
		if err != nil {
			return []entities.Student{}, err
		}

		for rows.Next() {
			var student entities.Student
			rows.Scan(&student.ID, &student.Name, &student.DOB, &student.Phone, &student.Branch, &student.Company.ID,
				&student.Status)
			students = append(students, student)
		}
	}

	if len(students) == 0 {
		return students, sql.ErrNoRows
	}

	return students, nil
}

// Create store to create a new student
func (store store) Create(student entities.Student) (entities.Student, error) {
	if err := store.db.Ping(); err != nil {
		return entities.Student{}, sql.ErrConnDone
	}

	if _, err := store.GetCompany(student.Company.ID); err != nil {
		return entities.Student{}, sql.ErrNoRows
	}

	student.ID = uuid.NewString()
	query := "INSERT INTO student (id, name, dob, phone, branch, companyID, status) VALUES(?, ?, ?, ?, ?, ?, ?)"
	_, err := store.db.Exec(query, student.ID, student.Name, student.DOB, student.Phone, student.Branch,
		student.Company.ID, student.Status)
	if err != nil {
		return entities.Student{}, err
	}
	return student, nil
}

// Update store to update a particular student
func (store store) Update(student entities.Student) (entities.Student, error) {
	if err := store.db.Ping(); err != nil {
		return entities.Student{}, sql.ErrConnDone
	}

	if _, err := store.GetById(student.ID); err != nil {
		return entities.Student{}, sql.ErrNoRows
	}

	if _, err := store.GetCompany(student.Company.ID); err != nil {
		return entities.Student{}, sql.ErrNoRows
	}

	query := "UPDATE student SET name=?, phone=?, dob=?, branch=?, companyID=?, status=? WHERE id=?"
	_, err := store.db.Exec(query, student.Name, student.Phone, student.DOB, student.Branch, student.Company.ID, student.Status,
		student.ID)

	if err != nil {
		return entities.Student{}, err
	}
	return student, nil
}

// Delete store to delete a particular student
func (store store) Delete(id string) error {
	if err := store.db.Ping(); err != nil {
		return sql.ErrConnDone
	}

	if _, err := store.GetById(id); err != nil {
		return sql.ErrNoRows
	}

	query := "DELETE FROM student WHERE id=?"
	_, err := store.db.Exec(query, id)
	if err != nil {
		return err
	}
	return nil
}

// GetCompany store to get company's detail a particular company id
func (store store) GetCompany(id string) (entities.Company, error) {
	if err := store.db.Ping(); err != nil {
		return entities.Company{}, sql.ErrConnDone
	}

	query := "SELECT * FROM company WHERE id=?"

	var company entities.Company
	row := store.db.QueryRow(query, id)
	err := row.Scan(&company.ID, &company.Name, &company.Category)

	if err != nil {
		return entities.Company{}, err
	}

	return company, nil
}
