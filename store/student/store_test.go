package student

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"student-placement-api/entities"
	"testing"
)

// createMockDb function to create sqlmock db
func createMockDB() (*sql.DB, sqlmock.Sqlmock, error) {
	return sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
}

// TestStore_Get to test function to test Student Get store
func TestStore_Get(t *testing.T) {
	db, mock, err := createMockDB()
	if err == nil {
		defer db.Close()
	}

	type inputStruct struct {
		name           string
		branch         string
		includeCompany bool
	}

	testcases := []struct {
		input    inputStruct
		expRes   []entities.Student
		mockRows []entities.Student
		expErr   error
		desc     string
	}{
		{
			inputStruct{"Test Company", "CSE", false},
			[]entities.Student{
				{
					"1", "Test Company", "12/12/2000", "CSE", "9876543210",
					entities.Company{}, "PENDING"},
			},
			[]entities.Student{
				{
					"1", "Test Company", "12/12/2000", "CSE", "9876543210",
					entities.Company{ID: "1", Name: "Test Company", Category: "MASS"}, "PENDING",
				},
			},
			nil, "Student with that ID exists so student should be returned",
		},
		{
			inputStruct{"Test Company", "CSE", true},
			[]entities.Student{
				{
					"1", "Test Company", "12/12/2000", "CSE", "9876543210",
					entities.Company{ID: "1", Name: "Test Company", Category: "MASS"}, "PENDING",
				},
			},
			[]entities.Student{
				{
					"1", "Test Company", "12/12/2000", "CSE", "9876543210",
					entities.Company{ID: "1", Name: "Test Company", Category: "MASS"}, "PENDING",
				},
			},
			nil, "Student with that ID exists so student should be returned",
		},
		{inputStruct{"Test Student", "CSE2", false}, []entities.Student{},
			[]entities.Student{
				{
					"1", "Test Company", "12/12/2000", "CSE", "9876543210",
					entities.Company{ID: "1", Name: "Test Company", Category: "MASS"}, "PENDING",
				},
			},
			sql.ErrNoRows, "Student with that name and branch doesn't exit"},
		{inputStruct{"Test Student 4", "CSE", false}, []entities.Student{},
			[]entities.Student{
				{
					"1", "Test Company", "12/12/2000", "CSE", "9876543210",
					entities.Company{ID: "1", Name: "Test Company", Category: "MASS"}, "PENDING",
				},
			},
			sql.ErrNoRows, "Student with that name and branch doesn't exit"},
		{inputStruct{"Test Student 4", "CSE2", false}, []entities.Student{},
			[]entities.Student{
				{
					"1", "Test Company", "12/12/2000", "CSE", "9876543210",
					entities.Company{ID: "1", Name: "Test Company", Category: "MASS"}, "PENDING",
				},
			},
			sql.ErrNoRows, "Student with that name and branch doesn't exit"},
	}

	for i, _ := range testcases {
		store := New(db)

		switch testcases[i].input.includeCompany {
		case true:
			rows := mock.NewRows([]string{"id", "name", "dob", "phone", "branch", "companyID", "companyName",
				"companyCategory", "status"})

			if testcases[i].expErr == nil {
				for j, _ := range testcases[i].expRes {
					rows.AddRow(testcases[i].mockRows[j].ID, testcases[i].mockRows[j].Name, testcases[i].mockRows[j].DOB,
						testcases[i].mockRows[j].Phone, testcases[i].mockRows[j].Branch,
						testcases[i].mockRows[j].Company.ID, testcases[i].mockRows[j].Company.Name,
						testcases[i].mockRows[j].Company.Category, testcases[i].mockRows[j].Status)
				}
			}

			mock.ExpectQuery("SELECT student.id AS id, student.name AS name, student.dob AS dob, "+
				"student.phone AS phone, student.branch AS branch, company.id AS companyID, company.name AS companyName, "+
				"company.category AS companyCategory, student.status AS status FROM student JOIN company "+
				"ON student.company_id=company.id WHERE student.name=? AND student.branch=?").
				WithArgs(testcases[i].input.name, testcases[i].input.branch).WillReturnRows(rows)

		case false:
			rows := mock.NewRows([]string{"id", "name", "dob", "phone", "branch", "status"})
			if testcases[i].expErr == nil {
				for j, _ := range testcases[i].expRes {
					rows.AddRow(testcases[i].mockRows[j].ID, testcases[i].mockRows[j].Name, testcases[i].mockRows[j].DOB,
						testcases[i].mockRows[j].Phone, testcases[i].mockRows[j].Branch,
						testcases[i].mockRows[j].Status)
				}
			}

			mock.ExpectQuery("SELECT id, name, dob, phone, branch, status FROM student WHERE student.name=? "+
				"AND student.branch=?").
				WithArgs(testcases[i].input.name, testcases[i].input.branch).WillReturnRows(rows)
		}

		actualRes, actualErr := store.Get(testcases[i].input.name, testcases[i].input.branch,
			testcases[i].input.includeCompany)

		if !reflect.DeepEqual(actualRes, testcases[i].expRes) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1, testcases[i].expRes,
				actualRes, testcases[i].desc)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expErr) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1, testcases[i].expErr,
				actualErr, testcases[i].desc)
		}
	}
}

// TestStore_GetById to test function to test Student GetById store
func TestStore_GetById(t *testing.T) {
	db, mock, err := createMockDB()
	if err == nil {
		defer db.Close()
	}

	testcases := []struct {
		id     string
		expRes entities.Student
		expErr error
		desc   string
	}{
		{"1", entities.Student{"1", "Test Student", "12/12/2000", "CSE",
			"9876543210", entities.Company{}, "PENDING"}, nil, "Student with that ID exists"},
		{"2", entities.Student{}, sql.ErrNoRows, "Student with that ID doesn't exist"},
	}

	for i, _ := range testcases {
		store := New(db)

		rows := mock.NewRows([]string{"id", "name", "dob", "phone", "branch", "status"})
		if testcases[i].expErr == nil {
			rows.AddRow(testcases[i].expRes.ID, testcases[i].expRes.Name, testcases[i].expRes.DOB,
				testcases[i].expRes.Phone, testcases[i].expRes.Branch, testcases[i].expRes.Status)
		}
		mock.ExpectQuery("SELECT student.id AS id, student.name AS name, student.dob AS dob, " +
			"student.phone AS phone, student.branch AS branch, student.status AS status " +
			"FROM student WHERE student.id=?").WithArgs(testcases[i].id).
			WillReturnRows(rows)

		actualRes, actualErr := store.GetById(testcases[i].id)

		if !reflect.DeepEqual(actualRes, testcases[i].expRes) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1, testcases[i].expRes,
				actualRes, testcases[i].desc)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expErr) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1, testcases[i].expErr,
				actualErr, testcases[i].desc)
		}
	}
}

// TestStore_Create to test function to test Student Create store
func TestStore_Create(t *testing.T) {
	db, mock, err := createMockDB()
	if err == nil {
		defer db.Close()
	}

	testcases := []struct {
		input  entities.Student
		expRes entities.Student
		expErr error
		desc   string
	}{
		{
			entities.Student{Name: "Test Student", DOB: "12/12/2000", Phone: "9876543210", Branch: "CSE",
				Company: entities.Company{ID: "1"}, Status: "PENDING"},
			entities.Student{ID: "1", Name: "Test Student", DOB: "12/12/2000", Phone: "9876543210",
				Branch: "CSE", Company: entities.Company{ID: "1"}, Status: "PENDING"},
			nil, "Student should be created",
		},
	}

	for i, _ := range testcases {
		store := New(db)

		mock.ExpectExec("INSERT INTO student (id, name, dob, phone, branch, company_id, status) "+
			"VALUES(?, ?, ?, ?, ?, ?, ?)").
			WithArgs(sqlmock.AnyArg(), testcases[i].input.Name, testcases[i].input.DOB,
				testcases[i].input.Phone, testcases[i].input.Branch, testcases[i].input.Company.ID,
				testcases[i].input.Status).
			WillReturnResult(sqlmock.NewResult(0, 1))

		_, actualErr := store.Create(testcases[i].input)

		if !reflect.DeepEqual(actualErr, testcases[i].expErr) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1, testcases[i].expErr,
				actualErr, testcases[i].desc)
		}
	}
}

// TestStore_Update to test function to test Student Update store
func TestStore_Update(t *testing.T) {
	db, mock, err := createMockDB()
	if err == nil {
		defer db.Close()
	}

	testcases := []struct {
		input  entities.Student
		expRes entities.Student
		expErr error
		desc   string
	}{
		{
			entities.Student{ID: "1", Name: "Test Student", DOB: "12/12/2000", Phone: "9876543210", Branch: "CSE",
				Company: entities.Company{ID: "1"}, Status: "PENDING"},
			entities.Student{ID: "1", Name: "Test Student", DOB: "12/12/2000", Phone: "9876543210",
				Branch: "CSE", Company: entities.Company{ID: "1"}, Status: "PENDING"},
			nil,
			"Student should be updated",
		},
	}

	for i, _ := range testcases {
		store := New(db)

		mock.ExpectExec("UPDATE student SET name=?, phone=?, dob=?, branch=?, company_id=?, status=? WHERE id=?").
			WithArgs(testcases[i].input.Name, testcases[i].input.Phone, testcases[i].input.DOB,
				testcases[i].input.Branch, testcases[i].input.Company.ID, testcases[i].input.Status,
				testcases[i].input.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		actualRes, actualErr := store.Update(testcases[i].input)

		if !reflect.DeepEqual(actualRes, testcases[i].expRes) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1, testcases[i].expRes,
				actualRes, testcases[i].desc)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expErr) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1, testcases[i].expErr,
				actualErr, testcases[i].desc)
		}
	}
}

// TestStore_Delete to test function to test Student Delete store
func TestStore_Delete(t *testing.T) {
	db, mock, err := createMockDB()
	if err == nil {
		defer db.Close()
	}

	testcases := []struct {
		id     string
		expErr error
		desc   string
	}{
		{"1", nil, "Student exists so should be deleted"},
	}

	for i, _ := range testcases {
		store := New(db)

		mock.ExpectExec("DELETE FROM student WHERE id=?").
			WithArgs(testcases[i].id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		actualErr := store.Delete(testcases[i].id)

		if !reflect.DeepEqual(actualErr, testcases[i].expErr) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1, testcases[i].expErr,
				actualErr, testcases[i].desc)
		}
	}
}

// TestStore_GetCompany to test function to test Student GetCompany store
func TestStore_GetCompany(t *testing.T) {
	db, mock, err := createMockDB()
	if err == nil {
		defer db.Close()
	}

	testcases := []struct {
		id     string
		expRes entities.Company
		expErr error
		desc   string
	}{
		{"1", entities.Company{"1", "Test Company", "MASS"}, nil,
			"Company with that ID exists"},
		{"1", entities.Company{}, sql.ErrNoRows,
			"Company with that ID doesn't exist"},
	}

	for i, _ := range testcases {
		store := New(db)

		rows := mock.NewRows([]string{"ID", "Name", "Category"})
		if testcases[i].expErr == nil {
			rows.AddRow(testcases[i].expRes.ID, testcases[i].expRes.Name, testcases[i].expRes.Category)
		}
		mock.ExpectQuery("SELECT * FROM company WHERE id=?").WithArgs(testcases[i].id).WillReturnRows(rows)

		actualRes, actualErr := store.GetCompany(testcases[i].id)

		if !reflect.DeepEqual(actualRes, testcases[i].expRes) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1, testcases[i].expRes,
				actualRes, testcases[i].desc)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expErr) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1, testcases[i].expErr,
				actualErr, testcases[i].desc)
		}
	}
}
