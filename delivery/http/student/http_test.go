package student

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"student-placement-api/entities"
	"testing"
)

const URL = "http://localhost:8080/student"

// TestHandler_Handler test function to test Student main handler function
func TestHandler_Handler(t *testing.T) {
	testcases := []struct {
		body          interface{}
		expecStatus   int
		expecResponse interface{}
		method        string
		description   string
	}{
		{
			"1",
			http.StatusOK,
			entities.Student{"1", "Test Student", "12/12/2000", "CSE",
				"9876543210", "1", "Pending"},
			http.MethodGet,
			"Student with that ID is present so a company should be returned and status code should be 200",
		},
		{
			entities.Student{
				Name:    "Test Student",
				DOB:     "12/12/2000",
				Phone:   "9876543210",
				Branch:  "CSE",
				Company: "1",
				Status:  "PENDING",
			},
			http.StatusCreated,
			nil,
			http.MethodPost,
			"Student should be added and status code should be 201",
		},
		{
			entities.Student{
				"1",
				"Test Student",
				"12/12/2000",
				"ECE",
				"9876543210",
				"1",
				"CORE",
			},
			http.StatusOK,
			entities.ResponseMessage{"Student Updated"},
			http.MethodPut,
			"Student should be updated and status code should be 200",
		},
		{
			"1",
			http.StatusOK,
			entities.ResponseMessage{"Company deleted"},
			http.MethodDelete,
			"Company with that ID should be deleted and status code should be 200",
		},
		{
			"1",
			http.StatusOK,
			entities.ResponseMessage{"Student deleted"},
			http.MethodDelete,
			"Student with that ID should be deleted and status code should be 200",
		},
	}

	for i := range testcases {
		var req *http.Request
		switch testcases[i].method {
		case http.MethodPost, http.MethodPut:
			reqBody, _ := json.Marshal(testcases[i].body)
			req = httptest.NewRequest(testcases[i].method, URL, bytes.NewBuffer(reqBody))
		default:
			req = httptest.NewRequest(testcases[i].method, URL+"?id="+testcases[i].body.(string), nil)
		}
		w := httptest.NewRecorder()
		handler := New(mockStudentService{})

		handler.Handler(w, req)

		if w.Code != testcases[i].expecStatus {
			t.Errorf("Test: %v\t Expected Code: %v\t Actual Code: %v\t Description: %v", i+1,
				testcases[i].expecStatus, w.Code, testcases[i].description)
		}
	}
}

// TestHandler_GetByID test function to test Student GetByID handler function
func TestHandler_Get(t *testing.T) {
	type bodyStruct struct {
		name           string
		branch         string
		includeCompany bool
	}

	testcases := []struct {
		body          bodyStruct
		expecStatus   int
		expecResponse interface{}
		description   string
	}{
		{
			bodyStruct{
				"Test Student",
				"CSE",
				false,
			},
			http.StatusOK,
			[]entities.Student{
				{"1", "Test Student", "12/12/2000", "CSE",
					"9876543210", "1", "Pending"},
			},
			"Student with that ID is present so a company should be returned and status code should be 200",
		},
		{
			bodyStruct{
				"Test Student",
				"CSE",
				true,
			},
			http.StatusOK,
			[]entities.Student{
				{"1", "Test Student", "12/12/2000", "CSE",
					"9876543210", entities.Company{"1", "Test Company", "MASS"}, "Pending"},
			},
			"Student with that ID is present so a company should be returned and status code should be 200",
		},
		{
			"2",
			http.StatusNotFound,
			entities.Student{},
			"Student with that ID is not present so empty json object should be returned wit status code 404",
		},
	}

	for i := range testcases {
		req := httptest.NewRequest(http.MethodGet, URL+"?id="+testcases[i].id, nil)
		w := httptest.NewRecorder()
		handler := New(mockStudentService{})

		handler.Get(w, req)

		if w.Code != testcases[i].expecStatus {
			t.Errorf("Test: %v\t Expected Code: %v\t Actual Code: %v\t Description: %v", i+1,
				testcases[i].expecStatus, w.Code, testcases[i].description)
		}
	}
}

// TestHandler_GetByID test function to test Student GetByID handler function
func TestHandler_GetByID(t *testing.T) {
	testcases := []struct {
		id            string
		expecStatus   int
		expecResponse entities.Student
		description   string
	}{
		{
			"1",
			http.StatusOK,
			entities.Student{"1", "Test Student", "12/12/2000", "CSE",
				"9876543210", "1", "Pending"},
			"Student with that ID is present so a company should be returned and status code should be 200",
		},
		{
			"2",
			http.StatusNotFound,
			entities.Student{},
			"Student with that ID is not present so empty json object should be returned wit status code 404",
		},
	}

	for i := range testcases {
		req := httptest.NewRequest(http.MethodGet, URL+"?id="+testcases[i].id, nil)
		w := httptest.NewRecorder()
		handler := New(mockStudentService{})

		handler.Get(w, req)

		if w.Code != testcases[i].expecStatus {
			t.Errorf("Test: %v\t Expected Code: %v\t Actual Code: %v\t Description: %v", i+1,
				testcases[i].expecStatus, w.Code, testcases[i].description)
		}
	}
}

// TestHandler_Create test function to test Student Create handler
func TestHandler_Create(t *testing.T) {
	testcases := []struct {
		body          entities.Student
		expecStatus   int
		expecResponse error
		description   string
	}{
		{
			entities.Student{
				Name:    "Test Student",
				DOB:     "12/12/2000",
				Phone:   "9876543210",
				Branch:  "CSE",
				Company: "1",
				Status:  "PENDING",
			},
			http.StatusCreated,
			nil,
			"Student should be added and status code should be 201",
		},
		{
			entities.Student{
				Name:    "Test Student 2",
				DOB:     "12/12/2000",
				Phone:   "9876543210",
				Branch:  "MCA",
				Company: "1",
				Status:  "PENDING",
			},
			http.StatusBadRequest,
			errors.New("invalid branch"),
			"Student should not be created as branch is not valid and status code should be 400",
		},
		{
			entities.Student{
				Name:    "Test Student 2",
				DOB:     "12/12/2000",
				Phone:   "98765432100000",
				Branch:  "CSE",
				Company: "1",
				Status:  "PENDING",
			},
			http.StatusBadRequest,
			errors.New("invalid phone"),
			"Student should not be created as phone is not valid and status code should be 400",
		},
		{
			entities.Student{
				Name:    "Test Student 3",
				DOB:     "12/12/2000",
				Phone:   "9876543210",
				Branch:  "CSE",
				Company: "1",
				Status:  "SUCCESS",
			},
			http.StatusBadRequest,
			errors.New("invalid status"),
			"Student should not be created as status is not valid and status code should be 400",
		},
	}

	for i := range testcases {
		reqBody, _ := json.Marshal(testcases[i].body)
		req := httptest.NewRequest(http.MethodPost, URL, bytes.NewReader(reqBody))
		w := httptest.NewRecorder()
		handler := New(mockStudentService{})

		handler.Create(w, req)

		if w.Code != testcases[i].expecStatus {
			t.Errorf("Test: %v\t Expected Code: %v\t Actual Code: %v\t Description: %v", i+1,
				testcases[i].expecStatus, w.Code, testcases[i].description)
		}
	}
}

// TestHandler_Update test function to test Student Update handler
func TestHandler_Update(t *testing.T) {
	testcases := []struct {
		body          entities.Student
		expecStatus   int
		expecResponse entities.ResponseMessage
		description   string
	}{
		{
			entities.Student{
				"1",
				"Test Student",
				"12/12/2000",
				"ECE",
				"9876543210",
				"1",
				"CORE",
			},
			http.StatusOK,
			entities.ResponseMessage{"Student Updated"},
			"Student should be updated and status code should be 200",
		},
		{
			entities.Student{
				"1",
				"Test Student",
				"12/12/2000",
				"ECE2",
				"9876543210",
				"1",
				"ACCEPTED",
			},
			http.StatusBadRequest,
			entities.ResponseMessage{"Invalid Branch"},
			"Student should not be update as branch is not valid and status code should be 400",
		},
		{
			entities.Student{
				"1",
				"Test Student",
				"12/12/2000",
				"ECE",
				"987654321013311",
				"1",
				"REJECTED",
			},
			http.StatusBadRequest,
			entities.ResponseMessage{"Invalid Phone"},
			"Student should not be update as phone no. is not valid and status code should be 400",
		},
		{
			entities.Student{
				"1",
				"Test Student",
				"12/12/2000",
				"ECE",
				"987654321013311",
				"1",
				"CORE",
			},
			http.StatusBadRequest,
			entities.ResponseMessage{"Invalid Status"},
			"Student should not be update as status is not valid and status code should be 400",
		},
		{
			entities.Student{
				"3",
				"Test Student",
				"12/12/2000",
				"ECE",
				"9876543210",
				"1",
				"PENDING",
			},
			http.StatusNotFound,
			entities.ResponseMessage{"No student with this ID"},
			"Student should not be update as no student with this id and status code should be 404",
		},
	}

	for i := range testcases {
		reqBody, _ := json.Marshal(testcases[i].body)
		req := httptest.NewRequest(http.MethodPut, URL, bytes.NewReader(reqBody))
		w := httptest.NewRecorder()
		handler := New(mockStudentService{})

		handler.Update(w, req)

		if w.Code != testcases[i].expecStatus {
			t.Errorf("Test: %v\t Expected Code: %v\t Actual Code: %v\t Description: %v", i+1,
				testcases[i].expecStatus, w.Code, testcases[i].description)
		}
	}
}

// TestHandler_Delete test function to test Student Delete handler function
func TestHandler_Delete(t *testing.T) {
	testcases := []struct {
		id            string
		expecStatus   int
		expecResponse entities.ResponseMessage
		description   string
	}{
		{
			"1",
			http.StatusOK,
			entities.ResponseMessage{"Student deleted"},
			"Student with that ID should be deleted and status code should be 200",
		},
		{
			"2",
			http.StatusNotFound,
			entities.ResponseMessage{"No student with that ID"},
			"Student with that ID is present so a company should be returned and status code should be 200",
		},
	}

	for i := range testcases {
		req := httptest.NewRequest(http.MethodDelete, URL+"?id="+testcases[i].id, nil)
		w := httptest.NewRecorder()
		handler := New(mockStudentService{})

		handler.Delete(w, req)

		if w.Code != testcases[i].expecStatus {
			t.Errorf("Test: %v\t Expected Code: %v\t Actual Code: %v\t Description: %v", i+1,
				testcases[i].expecStatus, w.Code, testcases[i].description)
		}
	}
}

type mockStudentService struct{}

// Get mock services for Get for Student
func (m mockStudentService) Get(name string, branch string, includeCompany bool) []entities.Student {
	return []entities.Student{}
}

// GetByID mock services for GetByID for Student
func (m mockStudentService) GetByID(id string) entities.Student {
	if id != "1" {
		return entities.Student{}
	}
	return entities.Student{"1", "Test Student", "12/12/2000", "CSE", "9876543210",
		"1", "Pending"}
}

// Create mock service for Create of Student
func (m mockStudentService) Create(student entities.Student) error {
	if student.Name == "" || student.Phone == "" || student.Company == "" ||
		student.Branch == "" || student.DOB == "" || student.Status == "" {
		return errors.New("all the fields are required, name, phone, dob, branch, company, status")
	}

	if len(student.Phone) < 10 || len(student.Phone) > 12 {
		return errors.New("invalid phone no.")
	}

	if !(student.Branch == "CSE" || student.Branch == "ISE" || student.Branch == "MECH" || student.Branch == "CIVIL" ||
		student.Branch == "ECE" || student.Branch == "EEE") {
		return errors.New("invalid branch")
	}

	if student.Status == "PENDING" || student.Status == "ACCEPTED" || student.Status == "REJECTED" {
		return errors.New("invalid status")
	}

	return nil
}

// Update mock service for Update of Student
func (m mockStudentService) Update(student entities.Student) error {
	if student.ID == "3" {
		return errors.New("student not found")
	}

	if student.ID == "" || student.Name == "" || student.Phone == "" || student.Company == "" ||
		student.Branch == "" || student.DOB == "" || student.Status == "" {
		return errors.New("all the fields are required, id, name, phone, dob, branch, company, status")
	}

	if len(student.Phone) < 10 || len(student.Phone) > 12 {
		return errors.New("invalid phone no.")
	}

	if !(student.Branch == "CSE" || student.Branch == "ISE" || student.Branch == "MECH" || student.Branch == "CIVIL" ||
		student.Branch == "ECE" || student.Branch == "EEE") {
		return errors.New("invalid branch")
	}

	if student.Status == "PENDING" || student.Status == "ACCEPTED" || student.Status == "REJECTED" {
		return errors.New("invalid status")
	}

	return nil
}

// Delete mock service for Delete of Student
func (m mockStudentService) Delete(id string) error {
	if id != "1" {
		return errors.New("student not found")
	}
	return nil
}
