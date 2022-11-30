package company

import (
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	_ "github.com/go-sql-driver/mysql"
	"reflect"
	"student-placement-api/entities"
	"testing"
)

// TestStore_GetByID to test function to test Company GetByID store
func TestStore_GetByID(t *testing.T) {
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
		{"2", entities.Company{}, errors.New("company not found")},
	}

	for i, _ := range testcases {
		store := New(db)

		rows := mock.NewRows([]string{"ID", "Name", "Category"})
		if testcases[i].expecErr == nil {
			rows.AddRow(testcases[i].expecRes.ID, testcases[i].expecRes.Name, testcases[i].expecRes.Category)
		}
		mock.ExpectQuery("SELECT * FROM company WHERE id=?").WithArgs(testcases[i].id).WillReturnRows(rows)

		actualRes, actualErr := store.GetByID(testcases[i].id)

		if !reflect.DeepEqual(actualRes, testcases[i].expecRes) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v", i+1, testcases[i].expecRes, actualRes)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expecErr) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v", i+1, testcases[i].expecErr, actualErr)
		}
	}
}

// TestStore_Create to test function to test Company Create store
func TestStore_Create(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err == nil {
		defer db.Close()
	}

	testcases := []struct {
		input    entities.Company
		expecRes entities.Company
		expecErr error
	}{
		{
			entities.Company{Name: "Test Company", Category: "MASS"},
			entities.Company{ID: "1", Name: "Test Company", Category: "MASS"},
			nil,
		},
		{
			entities.Company{Name: "Test Company 2", Category: "DREAM IT"},
			entities.Company{ID: "2", Name: "Test Company", Category: "DREAM IT"},
			nil,
		},
		{
			entities.Company{Name: "Test Company", Category: "OPEN DREAM"},
			entities.Company{ID: "1", Name: "Test Company", Category: "OPEN DREAM"},
			nil,
		},
	}

	for i, _ := range testcases {
		store := New(db)

		mock.ExpectExec("INSERT INTO company (id, name, category) VALUES(?, ?, ?)").
			WithArgs(sqlmock.AnyArg(), testcases[i].input.Name, testcases[i].input.Category).
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

// TestStore_Update to test function to test Company Update store
func TestStore_Update(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err == nil {
		defer db.Close()
	}

	testcases := []struct {
		input    entities.Company
		expecRes entities.Company
		expecErr error
	}{
		{
			entities.Company{ID: "1", Name: "Test Company", Category: "MASS"},
			entities.Company{ID: "1", Name: "Test Company", Category: "MASS"},
			nil,
		},
		{
			entities.Company{ID: "2", Name: "Test Company 2", Category: "DREAM IT"},
			entities.Company{},
			errors.New("company not found"),
		},
	}

	for i, _ := range testcases {
		store := New(db)

		mock.ExpectExec("UPDATE company SET name=?, category=? WHERE id=?").
			WithArgs(testcases[i].input.Name, testcases[i].input.Category, testcases[i].input.ID).
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

// TestStore_Delete to test function to test Company Delete store
func TestStore_Delete(t *testing.T) {
	db, mock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
	if err == nil {
		defer db.Close()
	}

	testcases := []struct {
		id       string
		expecRes entities.Company
		expecErr error
	}{
		{"1", entities.Company{ID: "1", Name: "Test Company", Category: "MASS"}, nil},
		{"2", entities.Company{}, errors.New("company not found")},
	}

	for i, _ := range testcases {
		store := New(db)

		mock.ExpectExec("DELETE FROM company WHERE id=?").
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
