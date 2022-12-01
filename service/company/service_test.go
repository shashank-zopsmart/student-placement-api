package company

import (
	"context"
	"database/sql"
	"errors"
	"reflect"
	"student-placement-api/entities"
	"testing"
)

// TestService_GetByID test function to test Company GetByID service function
func TestService_GetByID(t *testing.T) {
	testcases := []struct {
		id            string
		expecError    error
		expecResponse entities.Company
		description   string
	}{
		{
			"1",
			nil,
			entities.Company{"1", "Test Company", "MASS"},
			"Company with that ID is present so a company should be returned",
		},
		{
			"2",
			sql.ErrNoRows,
			entities.Company{},
			"Company with that ID is not present so error will be thrown",
		},
	}

	for i := range testcases {
		service := New(mockCompanyStore{})

		actualResponse, actualErr := service.GetByID(context.Background(), testcases[i].id)

		if !reflect.DeepEqual(actualResponse, testcases[i].expecResponse) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1,
				testcases[i].expecResponse, actualResponse, testcases[i].description)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expecError) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1,
				testcases[i].expecError, actualErr, testcases[i].description)
		}
	}
}

// TestService_Create test function to test Company Create service function
func TestService_Create(t *testing.T) {
	testcases := []struct {
		body          entities.Company
		expecError    error
		expecResponse entities.Company
		description   string
	}{
		{
			entities.Company{
				Name:     "Test Company",
				Category: "MASS",
			},
			nil,
			entities.Company{
				ID:       "1",
				Name:     "Test Company",
				Category: "MASS",
			},
			"Company should be added",
		},
		{
			entities.Company{
				Name:     "Test Company",
				Category: "NON IT",
			},
			errors.New("invalid category"),
			entities.Company{},
			"Company should not be added because of invalid category error will be thrown",
		},
	}

	for i := range testcases {
		service := New(mockCompanyStore{})

		actualResponse, actualErr := service.Create(context.Background(), testcases[i].body)

		if !reflect.DeepEqual(actualResponse, testcases[i].expecResponse) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1,
				testcases[i].expecResponse, actualResponse, testcases[i].description)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expecError) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1,
				testcases[i].expecError, actualErr, testcases[i].description)
		}
	}
}

// TestService_Update test function to test Company Update service function
func TestService_Update(t *testing.T) {
	testcases := []struct {
		body          entities.Company
		expecError    error
		expecResponse entities.Company
		description   string
	}{
		{
			entities.Company{
				"1",
				"Test Company",
				"CORE",
			},
			nil,
			entities.Company{
				"1",
				"Test Company",
				"CORE",
			},
			"Company should be updated",
		},
		{
			entities.Company{
				"1",
				"Test Company",
				"NON CORE",
			},
			errors.New("invalid category"),
			entities.Company{},
			"Company should not be updated because of invalid category error will be thrown",
		},
		{
			entities.Company{
				"2",
				"Test Company",
				"CORE",
			},
			sql.ErrNoRows,
			entities.Company{},
			"Company should not be updated because of invalid category error will be thrown",
		},
	}

	for i := range testcases {
		service := New(mockCompanyStore{})

		actualResponse, actualErr := service.Update(context.Background(), testcases[i].body)

		if !reflect.DeepEqual(actualResponse, testcases[i].expecResponse) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1,
				testcases[i].expecResponse, actualResponse, testcases[i].description)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expecError) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1,
				testcases[i].expecError, actualErr, testcases[i].description)
		}
	}
}

// TestService_Delete test function to test Company Delete service function
func TestService_Delete(t *testing.T) {
	testcases := []struct {
		id          string
		expecError  error
		description string
	}{
		{"1", nil, "Company with that ID should be deleted"},
		{
			"2", sql.ErrNoRows,
			"Company with that ID not is present error will be thrown",
		},
	}

	for i := range testcases {
		service := New(mockCompanyStore{})

		actualErr := service.Delete(context.Background(), testcases[i].id)

		if !reflect.DeepEqual(actualErr, testcases[i].expecError) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1,
				testcases[i].expecError, actualErr, testcases[i].description)
		}
	}
}

type mockCompanyStore struct{}

// Create mock store for Create of Company
func (m mockCompanyStore) Create(ctx context.Context, company entities.Company) (entities.Company, error) {
	switch company.Category {
	case "MASS", "DREAM IT", "OPEN DREAM", "CORE":
		return entities.Company{"1", "Test Company", "MASS"}, nil
	default:
		return entities.Company{}, errors.New("invalid category")
	}
}

// GetByID mock store for GetByID for Company
func (m mockCompanyStore) GetByID(ctx context.Context, id string) (entities.Company, error) {
	if id != "1" {
		return entities.Company{}, sql.ErrNoRows
	}
	return entities.Company{"1", "Test Company", "MASS"}, nil
}

// Update mock store for Update of Company
func (m mockCompanyStore) Update(ctx context.Context, company entities.Company) (entities.Company, error) {
	if company.ID != "1" {
		return entities.Company{}, sql.ErrNoRows
	}

	switch company.Category {
	case "MASS", "DREAM IT", "OPEN DREAM", "CORE":
		return entities.Company{"1", "Test Company", "CORE"}, nil
	default:
		return entities.Company{}, errors.New("invalid category")
	}
}

// Delete mock store for Delete of Company
func (m mockCompanyStore) Delete(ctx context.Context, id string) error {
	if id != "1" {
		return sql.ErrNoRows
	}
	return nil
}
