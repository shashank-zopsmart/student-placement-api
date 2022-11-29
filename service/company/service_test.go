package company

import (
	"errors"
	"reflect"
	"student-placement-api/entities"
	"testing"
)

// TestService_Get test function to test Company Get service function
func TestService_Get(t *testing.T) {
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
			"Company with that ID is present so a company should be returned and status code should be 200",
		},
		{
			"2",
			errors.New("company not found"),
			entities.Company{},
			"Company with that ID is not present so error will be thrown",
		},
	}

	for i := range testcases {
		service := New(mockCompanyStore{})

		actualResponse, actualErr := service.GetByID(testcases[i].id)

		if !reflect.DeepEqual(actualResponse, testcases[i].expecResponse) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1,
				testcases[i].expecResponse, actualResponse, testcases[i].description)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expecError) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1,
				testcases[i].expecResponse, actualResponse, testcases[i].description)
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
			"Company should be added and status code should be 201",
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

		actualResponse, actualErr := service.Create(testcases[i].body)

		if !reflect.DeepEqual(actualResponse, testcases[i].expecResponse) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1,
				testcases[i].expecResponse, actualResponse, testcases[i].description)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expecError) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1,
				testcases[i].expecResponse, actualResponse, testcases[i].description)
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
			"Company should be updated and status code should be 200",
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
			errors.New("company not found"),
			entities.Company{},
			"Company should not be updated because of invalid category error will be thrown",
		},
	}

	for i := range testcases {
		service := New(mockCompanyStore{})

		actualResponse, actualErr := service.Update(testcases[i].body)

		if !reflect.DeepEqual(actualResponse, testcases[i].expecResponse) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1,
				testcases[i].expecResponse, actualResponse, testcases[i].description)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expecError) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1,
				testcases[i].expecResponse, actualResponse, testcases[i].description)
		}
	}
}

// TestService_Delete test function to test Company Delete service function
func TestService_Delete(t *testing.T) {
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
			"Company with that ID should be deleted and status code should be 200",
		},
		{
			"2",
			errors.New("company not found"),
			entities.Company{},
			"Company with that ID not is present error will be thrown",
		},
	}

	for i := range testcases {
		service := New(mockCompanyStore{})

		actualResponse, actualErr := service.GetByID(testcases[i].id)

		if !reflect.DeepEqual(actualResponse, testcases[i].expecResponse) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1,
				testcases[i].expecResponse, actualResponse, testcases[i].description)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expecError) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1,
				testcases[i].expecResponse, actualResponse, testcases[i].description)
		}
	}
}

type mockCompanyStore struct{}

// GetByID mock store for GetByID for Company
func (m mockCompanyStore) GetByID(id string) (entities.Company, error) {
	if id != "1" {
		return entities.Company{}, errors.New("company not found")
	}
	return entities.Company{"1", "Test Company", "MASS"}, nil
}

// Create mock store for Create of Company
func (m mockCompanyStore) Create(company entities.Company) (entities.Company, error) {
	switch company.Category {
	case "MASS", "DREAM IT", "OPEN DREAM", "CORE":
		return entities.Company{"1", "Test Company", "MASS"}, nil
	default:
		return entities.Company{}, errors.New("invalid category")
	}
}

// Update mock store for Update of Company
func (m mockCompanyStore) Update(company entities.Company) (entities.Company, error) {
	if company.ID == "3" {
		return entities.Company{}, errors.New("company not found")
	}

	switch company.Category {
	case "MASS", "DREAM IT", "OPEN DREAM", "CORE":
		return entities.Company{}, nil
	default:
		return entities.Company{}, errors.New("invalid category")
	}
}

// Delete mock store for Delete of Company
func (m mockCompanyStore) Delete(id string) (entities.Company, error) {
	if id != "1" {
		return entities.Company{}, errors.New("company not found")
	}
	return entities.Company{}, nil
}
