package student

import (
	"context"
	"student-placement-api/errors"

	"fmt"
	"reflect"
	"student-placement-api/entities"
	"testing"
)

// TestService_Get test function to test Student Get service function
func TestService_Get(t *testing.T) {
	testcases := []struct {
		name           string
		branch         string
		includeCompany bool
		expecError     error
		expecResponse  []entities.Student
		description    string
	}{
		{
			"Student",
			"CSE",
			false,
			nil,
			[]entities.Student{
				{"1", "Student", "12/12/2000", "CSE",
					"9876543210", entities.Company{ID: "1"}, "Pending"},
			},
			"Student with that name and branch is present so a student should be returned",
		},
		{
			"Student",
			"CSE",
			true,
			nil,
			[]entities.Student{
				{"1", "Student", "12/12/2000", "CSE",
					"9876543210", entities.Company{"1", "Test Company", "MASS"}, "Pending"},
			},
			"Student with that name and branch is present includeCompany flag is true so a student " +
				"with company detail should be returned",
		},
		{
			"Student",
			"CSE2",
			false,
			errors.EntityNotFound{"Student"},
			[]entities.Student{},
			"Student with that name and branch branch is not present error will be thrown",
		},
		{
			"Student5",
			"CSE",
			false,
			errors.EntityNotFound{"Student"},
			[]entities.Student{},
			"Student with that name and branch branch is not present error will be thrown",
		},
	}

	for i := range testcases {
		service := New(mockStudentStore{})

		actualResponse, actualErr := service.Get(context.Background(), testcases[i].name, testcases[i].branch, testcases[i].includeCompany)

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

// TestService_GetByID test function to test Student GetByID service function
func TestService_GetByID(t *testing.T) {
	testcases := []struct {
		id            string
		expecError    error
		expecResponse entities.Student
		description   string
	}{
		{
			"1",
			nil,
			entities.Student{"1", "Test Student", "12/12/2000", "CSE",
				"9876543210", entities.Company{ID: "1"}, "PENDING"},
			"Student with that ID is present so a student should be returned",
		},
		{
			"2",
			errors.EntityNotFound{"Student"},
			entities.Student{},
			"Student with that ID is not present so error will be returned",
		},
	}

	for i := range testcases {
		service := New(mockStudentStore{})

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

// TestService_Create test function to test Student Create service
func TestService_Create(t *testing.T) {
	testcases := []struct {
		body          entities.Student
		expecError    error
		expecResponse entities.Student
		description   string
	}{
		{
			entities.Student{
				Name:    "Test Student",
				DOB:     "12/12/2000",
				Phone:   "9876543210",
				Branch:  "CSE",
				Company: entities.Company{ID: "1"},
				Status:  "PENDING",
			},
			nil,
			entities.Student{
				ID:      "1",
				Name:    "Test Student",
				DOB:     "12/12/2000",
				Phone:   "9876543210",
				Branch:  "CSE",
				Company: entities.Company{ID: "1"},
				Status:  "PENDING",
			},
			"Student should be added",
		},
		{
			entities.Student{
				Name:    "Test Student",
				DOB:     "12/12/2010",
				Phone:   "9876543210",
				Branch:  "CSE",
				Company: entities.Company{ID: "1"},
				Status:  "PENDING",
			},
			errors.InvalidParams{"Student doesn't meet minimum age requirement"},
			entities.Student{},
			"Student should be added as doesn't meet minimum age requirement",
		},
		{
			entities.Student{
				Name:    "Test Student 2",
				DOB:     "12/12/2000",
				Phone:   "9876543210",
				Branch:  "MCA",
				Company: entities.Company{ID: "1"},
				Status:  "PENDING",
			},
			errors.InvalidParams{"Invalid Branch"},
			entities.Student{},
			"Student should not be created as branch is not valid",
		},
		{
			entities.Student{
				Name:    "Test Student 2",
				DOB:     "12/12/2000",
				Phone:   "98765432100000",
				Branch:  "CSE",
				Company: entities.Company{ID: "1"},
				Status:  "PENDING",
			},
			errors.InvalidParams{"Phone must be of 10-12 digits"},
			entities.Student{},
			"Student should not be created as phone is not valid",
		},
		{
			entities.Student{
				Name:    "Test Student 3",
				DOB:     "12/12/2000",
				Phone:   "9876543210",
				Branch:  "CSE",
				Company: entities.Company{ID: "1"},
				Status:  "SUCCESS",
			},
			errors.InvalidParams{"Invalid Status"},
			entities.Student{},
			"Student should not be created as status is not valid",
		},
	}

	for i := range testcases {
		service := New(mockStudentStore{})

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

// TestService_Update test function to test Student Update service
func TestService_Update(t *testing.T) {
	testcases := []struct {
		body          entities.Student
		expecError    error
		expecResponse entities.Student
		description   string
	}{
		{
			entities.Student{
				"1",
				"Test Student",
				"12/12/2000",
				"ECE",
				"9876543210",
				entities.Company{ID: "1"},
				"PENDING",
			},
			nil,
			entities.Student{"1",
				"Test Student",
				"12/12/2000",
				"ECE",
				"9876543210",
				entities.Company{ID: "1"},
				"PENDING",
			},
			"Student should be updated",
		},
		{
			entities.Student{
				"1",
				"Test Student",
				"12/12/2010",
				"ECE",
				"9876543210",
				entities.Company{ID: "1"},
				"PENDING",
			},
			errors.InvalidParams{"Student doesn't meet minimum age requirement"},
			entities.Student{},
			"Student should be added as doesn't meet minimum age requirement",
		},
		{
			entities.Student{
				"1",
				"Test Student",
				"12/12/2000",
				"ECE2",
				"9876543210",
				entities.Company{ID: "1"},
				"ACCEPTED",
			},
			errors.InvalidParams{"Invalid Branch"},
			entities.Student{},
			"Student should not be update as branch is not valid",
		},
		{
			entities.Student{
				"1",
				"Test Student",
				"12/12/2000",
				"ECE",
				"987654321013311",
				entities.Company{ID: "1"},
				"REJECTED",
			},
			errors.InvalidParams{"Phone must be of 10-12 digits"},
			entities.Student{},
			"Student should not be update as phone no. is not valid",
		},
		{
			entities.Student{
				"1",
				"Test Student",
				"12/12/2000",
				"ECE",
				"98713311",
				entities.Company{ID: "1"},
				"CORE",
			},
			errors.InvalidParams{"Phone must be of 10-12 digits"},
			entities.Student{},
			"Student should not be update as status is not valid",
		},
		{
			entities.Student{
				"1",
				"Test Student",
				"12/12/2000",
				"ECE",
				"987654321a",
				entities.Company{ID: "1"},
				"CORE",
			},
			errors.InvalidParams{"Invalid Phone"},
			entities.Student{},
			"Student should not be update as status is not valid",
		},
		{
			entities.Student{
				"3",
				"Test Student",
				"12/12/2000",
				"ECE",
				"9876543210",
				entities.Company{ID: "1"},
				"PENDING",
			},
			errors.EntityNotFound{"Student"},
			entities.Student{},
			"Student should not be update as no student with this id",
		},
	}

	for i := range testcases {
		service := New(mockStudentStore{})

		actualResponse, actualErr := service.Update(context.Background(), testcases[i].body)

		if !reflect.DeepEqual(actualResponse, testcases[i].expecResponse) {
			fmt.Println(actualErr)
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1,
				testcases[i].expecResponse, actualResponse, testcases[i].description)
		}

		if !reflect.DeepEqual(actualErr, testcases[i].expecError) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1,
				testcases[i].expecError, actualErr, testcases[i].description)
		}
	}
}

// TestService_Delete test function to test Student Delete service function
func TestService_Delete(t *testing.T) {
	testcases := []struct {
		id          string
		expecError  error
		description string
	}{
		{"1", nil, "Student with that ID should be deleted"},
		{
			"2", errors.EntityNotFound{"Student"},
			"Student with that ID is not present so error will be thrown",
		},
	}

	for i := range testcases {
		service := New(mockStudentStore{})

		actualErr := service.Delete(context.Background(), testcases[i].id)

		if !reflect.DeepEqual(actualErr, testcases[i].expecError) {
			t.Errorf(" Test: %v\n Expected: %v\n Actual: %v\n Description: %v", i+1,
				testcases[i].expecError, actualErr, testcases[i].description)
		}
	}
}

type mockStudentStore struct{}

// Get mock store for Get for Student
func (m mockStudentStore) Get(ctx context.Context, name string, branch string, includeCompany bool) ([]entities.Student, error) {
	if name == "Student" && branch == "CSE" {
		if includeCompany == true {
			return []entities.Student{
				{"1", "Student", "12/12/2000", "CSE",
					"9876543210", entities.Company{"1", "Test Company", "MASS"}, "Pending"},
			}, nil
		}
		return []entities.Student{
			{"1", "Student", "12/12/2000", "CSE",
				"9876543210", entities.Company{ID: "1"}, "Pending"},
		}, nil
	}
	return []entities.Student{}, errors.EntityNotFound{"Student"}
}

// GetById mock store for GetById for Student
func (m mockStudentStore) GetById(ctx context.Context, id string) (entities.Student, error) {
	if id != "1" {
		return entities.Student{}, errors.EntityNotFound{"Student"}
	}
	return entities.Student{"1", "Test Student", "12/12/2000", "CSE", "9876543210",
		entities.Company{ID: "1"}, "PENDING"}, nil
}

// Create mock store for Create of Student
func (m mockStudentStore) Create(ctx context.Context, student entities.Student) (entities.Student, error) {
	return entities.Student{
		ID:      "1",
		Name:    "Test Student",
		DOB:     "12/12/2000",
		Phone:   "9876543210",
		Branch:  "CSE",
		Company: entities.Company{ID: "1"},
		Status:  "PENDING",
	}, nil
}

// Update mock store for Update of Student
func (m mockStudentStore) Update(ctx context.Context, student entities.Student) (entities.Student, error) {
	return entities.Student{"1", "Test Student", "12/12/2000", "ECE", "9876543210",
		entities.Company{ID: "1"}, "PENDING"}, nil
}

// Delete mock store for Delete of Student
func (m mockStudentStore) Delete(ctx context.Context, id string) error {
	return nil
}

// GetCompany mock store for GetCompany of Student
func (m mockStudentStore) GetCompany(ctx context.Context, id string) (entities.Company, error) {
	return entities.Company{"1", "Test Company", "MASS"}, nil
}
