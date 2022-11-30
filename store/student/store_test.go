package student

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"reflect"
	"student-placement-api/entities"
	"testing"
)

// TestStore_Get to test function to test Student Get store
func TestStore_Get(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
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
		expecRes []entities.Student
		expecErr error
	}{
		{
			inputStruct{"Test Company", "CSE", false},
			[]entities.Student{
				{
					"1", "Test Company", "12/12/2000", "CSE", "9876543210",
					entities.Company{ID: "1"}, "PENDING"},
			},
			nil,
		},
		{
			inputStruct{"Test Company", "CSE", true},
			[]entities.Student{
				{
					"1", "Test Company", "12/12/2000", "CSE", "9876543210",
					entities.Company{ID: "1", Name: "Test Company", Category: "MASS"}, "PENDING",
				},
			},
			nil,
		},
		{inputStruct{"Test Company", "CSE2", false}, []entities.Student{},
			errors.New("student not found")},
		{inputStruct{"Test Company 4", "CSE", false}, []entities.Student{},
			errors.New("student not found")},
		{inputStruct{"Test Company3", "CSE2", false}, []entities.Student{},
			errors.New("student not found")},
	}

	for i, _ := range testcases {
		store := New(db)

		switch testcases[i].input.includeCompany {
		case true:
			rows := mock.NewRows([]string{"id", "name", "dob", "phone", "branch", "companyID", "companyName",
				"companyCategory", "status"})

			if testcases[i].expecErr == nil {
				for j, _ := range testcases[i].expecRes {
					rows.AddRow(testcases[i].expecRes[j].ID, testcases[i].expecRes[j].Name, testcases[i].expecRes[j].DOB,
						testcases[i].expecRes[j].Phone, testcases[i].expecRes[j].Branch,
						testcases[i].expecRes[j].Company.ID, testcases[i].expecRes[j].Company.Name,
						testcases[i].expecRes[j].Company.Category, testcases[i].expecRes[j].Status)
				}
			}

			mock.ExpectQuery("SELECT student.id AS id, student.name AS name, student.dob AS dob, "+
				"student.phone AS phone, student.branch AS branch, company.id AS companyID, company.name AS companyName, "+
				"company.category AS companyCategory, student.status AS status FROM student JOIN company "+
				"ON student.id=company.id WHERE student.name=? AND student.branch=?").
				WithArgs(testcases[i].input.name, testcases[i].input.branch).WillReturnRows(rows)

		case false:
			rows := mock.NewRows([]string{"id", "name", "dob", "phone", "branch", "companyID", "status"})
			if testcases[i].expecErr == nil {
				for j, _ := range testcases[i].expecRes {
					rows.AddRow(testcases[i].expecRes[j].ID, testcases[i].expecRes[j].Name, testcases[i].expecRes[j].DOB,
						testcases[i].expecRes[j].Phone, testcases[i].expecRes[j].Branch,
						testcases[i].expecRes[j].Company.ID, testcases[i].expecRes[j].Status)
				}
			}

			mock.ExpectQuery("SELECT * FROM student WHERE student.name=? AND student.branch=?").
				WithArgs(testcases[i].input.name, testcases[i].input.branch).WillReturnRows(rows)
		}

		actualRes, actualErr := store.Get(testcases[i].input.name, testcases[i].input.branch,
			testcases[i].input.includeCompany)

		if !reflect.DeepEqual(actualRes, testcases[i].expecRes) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v", i+1, testcases[i].expecRes, actualRes)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expecErr) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v", i+1, testcases[i].expecErr, actualErr)
		}
	}
}

// TestStore_GetById to test function to test Student GetById store
func TestStore_GetById(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err == nil {
		defer db.Close()
	}

	testcases := []struct {
		id       string
		expecRes entities.Student
		expecErr error
	}{
		{"1", entities.Student{"1", "Test Student", "12/12/2000", "CSE",
			"9876543210", entities.Company{ID: "1"}, "PENDING"}, nil},
		{"2", entities.Student{}, errors.New("student not found")},
	}

	for i, _ := range testcases {
		store := New(db)

		rows := mock.NewRows([]string{"id", "name", "dob", "phone", "branch", "companyID", "status"})
		if testcases[i].expecErr == nil {
			rows.AddRow(testcases[i].expecRes.ID, testcases[i].expecRes.Name, testcases[i].expecRes.DOB,
				testcases[i].expecRes.Phone, testcases[i].expecRes.Branch, testcases[i].expecRes.Company.ID,
				testcases[i].expecRes.Status)
		}
		mock.ExpectQuery("SELECT student.id AS id, student.name AS name, student.dob AS dob, " +
			"student.phone AS phone, student.branch AS branch, company.id AS companyID, student.status AS status " +
			"FROM student JOIN company ON student.id=company.id WHERE student.id=?").WithArgs(testcases[i].id).
			WillReturnRows(rows)

		actualRes, actualErr := store.GetById(testcases[i].id)

		if !reflect.DeepEqual(actualRes, testcases[i].expecRes) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v", i+1, testcases[i].expecRes, actualRes)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expecErr) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v", i+1, testcases[i].expecErr, actualErr)
		}
	}
}

// TestStore_Create to test function to test Student Create store
func TestStore_Create(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err == nil {
		defer db.Close()
	}

	testcases := []struct {
		input    entities.Student
		expecRes entities.Student
		expecErr error
	}{
		{
			entities.Student{Name: "Test Student", DOB: "12/12/2000", Phone: "9876543210", Branch: "CSE",
				Company: entities.Company{ID: "1"}, Status: "PENDING"},
			entities.Student{ID: "1", Name: "Test Student", DOB: "12/12/2000", Phone: "9876543210",
				Branch: "CSE", Company: entities.Company{ID: "1"}, Status: "PENDING"},
			nil,
		},
	}

	for i, _ := range testcases {
		store := New(db)

		mock.ExpectExec("INSERT INTO student (id, name, dob, phone, branch, companyID, status) "+
			"VALUES(?, ?, ?, ?, ?, ?, ?)").
			WithArgs(sqlmock.AnyArg(), testcases[i].input.Name, testcases[i].input.DOB, testcases[i].input.DOB,
				testcases[i].input.Phone, testcases[i].input.Branch, testcases[i].input.Company.ID,
				testcases[i].input.Status).
			WillReturnResult(sqlmock.NewResult(0, 1))

		actualRes, actualErr := store.Create(testcases[i].input)

		if !reflect.DeepEqual(actualRes, testcases[i].expecRes) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v", i+1, testcases[i].expecRes, actualRes)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expecErr) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v", i+1, testcases[i].expecErr, actualErr)
		}
	}
}

// TestStore_Update to test function to test Student Update store
func TestStore_Update(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err == nil {
		defer db.Close()
	}

	testcases := []struct {
		input    entities.Student
		expecRes entities.Student
		expecErr error
	}{
		{
			entities.Student{ID: "1", Name: "Test Student", DOB: "12/12/2000", Phone: "9876543210", Branch: "CSE",
				Company: entities.Company{ID: "1"}, Status: "PENDING"},
			entities.Student{ID: "1", Name: "Test Student", DOB: "12/12/2000", Phone: "9876543210",
				Branch: "CSE", Company: entities.Company{ID: "1"}, Status: "PENDING"},
			nil,
		},
		{
			entities.Student{ID: "1", Name: "Test Student", DOB: "12/12/2000", Phone: "9876543210", Branch: "CSE",
				Company: entities.Company{ID: "2"}, Status: "PENDING"},
			entities.Student{ID: "1", Name: "Test Student", DOB: "12/12/2000", Phone: "9876543210",
				Branch: "CSE", Company: entities.Company{ID: "2"}, Status: "PENDING"},
			nil,
		},
		{
			entities.Student{ID: "1", Name: "Test Student", DOB: "12/12/2000", Phone: "9876543210", Branch: "CSE",
				Company: entities.Company{ID: "3"}, Status: "PENDING"},
			entities.Student{},
			errors.New(""),
		},
	}

	for i, _ := range testcases {
		store := New(db)

		mock.ExpectExec("UPDATE student SET name=?, phone=?, dob=?, branch=?, company=?, status=? WHERE id=?").
			WithArgs(testcases[i].input.Name, testcases[i].input.DOB, testcases[i].input.DOB,
				testcases[i].input.Phone, testcases[i].input.Branch, testcases[i].input.Company.ID,
				testcases[i].input.Status, testcases[i].input.ID).
			WillReturnResult(sqlmock.NewResult(0, 1))

		actualRes, actualErr := store.Update(testcases[i].input)

		if !reflect.DeepEqual(actualRes, testcases[i].expecRes) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v", i+1, testcases[i].expecRes, actualRes)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expecErr) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v", i+1, testcases[i].expecErr, actualErr)
		}
	}
}

// TestStore_Delete to test function to test Student Delete store
func TestStore_Delete(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err == nil {
		defer db.Close()
	}

	testcases := []struct {
		id       string
		expecRes entities.Student
		expecErr error
	}{
		{"1", entities.Student{"1", "Test Student", "12/12/2000", "CSE",
			"9876543210", entities.Company{ID: "1"}, "PENDING"}, nil},
		{"2", entities.Student{}, errors.New("student not found")},
	}

	for i, _ := range testcases {
		store := New(db)

		mock.ExpectExec("DELETE FROM student WHERE id=?").
			WithArgs(testcases[i].id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		actualRes, actualErr := store.Delete(testcases[i].id)

		if !reflect.DeepEqual(actualRes, testcases[i].expecRes) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v", i+1, testcases[i].expecRes, actualRes)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expecErr) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v", i+1, testcases[i].expecErr, actualErr)
		}
	}
}

// TestStore_GetCompany to test function to test Student GetCompany store
func TestStore_GetCompany(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err == nil {
		defer db.Close()
	}

	testcases := []struct {
		id       string
		expecRes entities.Company
		expecErr error
	}{
		{"1", entities.Company{"1", "Test Company", "MASS"}, nil},
		{"1", entities.Company{}, errors.New("company not found")},
	}

	for i, _ := range testcases {
		store := New(db)

		rows := mock.NewRows([]string{"ID", "Name", "Category"})
		if testcases[i].expecErr == nil {
			rows.AddRow(testcases[i].expecRes.ID, testcases[i].expecRes.Name, testcases[i].expecRes.Category)
		}
		mock.ExpectQuery("SELECT * FROM company WHERE id=?").WithArgs(testcases[i].id).WillReturnRows(rows)

		actualRes, actualErr := store.GetCompany(testcases[i].id)

		if !reflect.DeepEqual(actualRes, testcases[i].expecRes) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v", i+1, testcases[i].expecRes, actualRes)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expecErr) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v", i+1, testcases[i].expecErr, actualErr)
		}
	}
}
