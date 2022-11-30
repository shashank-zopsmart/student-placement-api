package company

import (
	"database/sql"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"student-placement-api/entities"
	"testing"
)

// createMockDb function to create sqlmock db
func createMockDB() (*sql.DB, sqlmock.Sqlmock, error) {
	return sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
}

// TestStore_GetByID to test function to test Company GetByID store
func TestStore_GetByID(t *testing.T) {
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
		{"2", entities.Company{}, sql.ErrNoRows, "No company exists with that ID"},
		{"2", entities.Company{}, sql.ErrConnDone, "Database connection is closed"},
	}

	for i, _ := range testcases {
		store := New(db)

		rows := mock.NewRows([]string{"ID", "Name", "Category"})
		if testcases[i].expErr == nil {
			rows.AddRow(testcases[i].expRes.ID, testcases[i].expRes.Name, testcases[i].expRes.Category)
		}

		mock.ExpectQuery("SELECT * FROM company WHERE id=?").WithArgs(testcases[i].id).WillReturnRows(rows)

		actualRes, actualErr := store.GetByID(testcases[i].id)

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

// TestStore_Create to test function to test Company Create store
func TestStore_Create(t *testing.T) {
	db, mock, err := createMockDB()
	if err == nil {
		defer db.Close()
	}

	testcases := []struct {
		input  entities.Company
		expRes entities.Company
		expErr error
		desc   string
	}{
		{
			entities.Company{Name: "Test Company", Category: "MASS"},
			entities.Company{ID: "1", Name: "Test Company", Category: "MASS"}, nil,
			"Company should be inserted",
		},
		{
			entities.Company{Name: "Test Company", Category: "OPEN DREAM"},
			entities.Company{ID: "1", Name: "Test Company", Category: "OPEN DREAM"}, sql.ErrConnDone,
			"Database connection closed",
		},
	}

	for i, _ := range testcases {
		store := New(db)

		mock.ExpectExec("INSERT INTO company (id, name, category) VALUES(?, ?, ?)").
			WithArgs("1", testcases[i].input.Name, testcases[i].input.Category).
			WillReturnResult(sqlmock.NewResult(0, 1))

		actualRes, actualErr := store.Create(testcases[i].input)

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

// TestStore_Update to test function to test Company Update store
func TestStore_Update(t *testing.T) {
	db, mock, err := createMockDB()
	if err == nil {
		defer db.Close()
	}

	testcases := []struct {
		input  entities.Company
		expRes entities.Company
		expErr error
		desc   string
	}{
		{
			entities.Company{ID: "1", Name: "Test Company", Category: "MASS"},
			entities.Company{ID: "1", Name: "Test Company", Category: "MASS"}, nil,
			"Company detail should be updated",
		},
		{
			entities.Company{ID: "2", Name: "Test Company 2", Category: "DREAM IT"},
			entities.Company{}, sql.ErrNoRows,
			"Company with that ID doesn't exist",
		},
		{
			entities.Company{ID: "1", Name: "Test Company", Category: "MASS"},
			entities.Company{}, sql.ErrConnDone,
			"Database connection closed",
		},
	}

	for i, _ := range testcases {
		store := New(db)

		mock.ExpectExec("UPDATE company SET name=?, category=? WHERE id=?").
			WithArgs(testcases[i].input.Name, testcases[i].input.Category, testcases[i].input.ID).
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

// TestStore_Delete to test function to test Company Delete store
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
		{"1", nil, "Company should be deleted"},
		{"2", sql.ErrNoRows, "Company with that ID doesn't exist"},
		{"2", sql.ErrConnDone, "Database connection closed"},
	}

	for i, _ := range testcases {
		store := New(db)

		mock.ExpectExec("DELETE FROM company WHERE id=?").
			WithArgs(testcases[i].id).
			WillReturnResult(sqlmock.NewResult(0, 1))

		actualErr := store.Delete(testcases[i].id)

		if !reflect.DeepEqual(actualErr, testcases[i].expErr) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1, testcases[i].expErr,
				actualErr, testcases[i].desc)
		}
	}
}
