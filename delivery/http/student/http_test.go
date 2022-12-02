package student

import (
	"bytes"
	"context"
	"encoding/json"
	"student-placement-api/errors"

	"fmt"
	"net/http"
	"net/http/httptest"
	"student-placement-api/entities"
	"testing"
)

const URL = "http://localhost:8080/student"

// TestHandler_Handler test function to test Student main handler function
func TestHandler_Handler(t *testing.T) {
	type bodyStruct struct {
		id             string
		name           string
		branch         string
		includeCompany string
	}
	testcases := []struct {
		body          interface{}
		expecStatus   int
		expecResponse interface{}
		method        string
		description   string
	}{
		{
			bodyStruct{id: "1"},
			http.StatusOK,
			entities.Student{"1", "Test Student", "12/12/2000", "CSE",
				"9876543210", entities.Company{ID: "1"}, "Pending"},
			http.MethodGet,
			"Student with that ID is present so a company should be returned and status code should be 200",
		},
		{
			bodyStruct{name: "Student", branch: "CSE", includeCompany: "false"},
			http.StatusOK,
			entities.Student{"1", "Student", "12/12/2000", "CSE",
				"9876543210", entities.Company{}, "Pending"},
			http.MethodGet,
			"Student with that ID is present so a company should be returned and status code should be 200",
		},
		{
			bodyStruct{"", "", "", ""},
			http.StatusBadRequest,
			entities.Student{},
			http.MethodGet,
			"No id, name, branch and includeCompany given status code should be 400",
		},
		{
			entities.Student{
				Name:    "Test Student",
				DOB:     "12/12/2000",
				Phone:   "9876543210",
				Branch:  "CSE",
				Company: entities.Company{ID: "1"},
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
				entities.Company{ID: "1"},
				"CORE",
			},
			http.StatusBadRequest,
			entities.ResponseMessage{"Error: invalid branch"},
			http.MethodPut,
			"Student should be updated and status code should be 200",
		},
		{
			"1",
			http.StatusOK,
			entities.ResponseMessage{"Student deleted"},
			http.MethodDelete,
			"Student with that ID should be deleted and status code should be 200",
		},
		{
			"1",
			http.StatusOK,
			entities.ResponseMessage{"Student deleted"},
			http.MethodDelete,
			"Student with that ID should be deleted and status code should be 200",
		},
		{
			entities.Student{
				"1",
				"Test Student",
				"12/12/2000",
				"ECE",
				"9876543210",
				entities.Company{ID: "1"},
				"CORE",
			},
			http.StatusMethodNotAllowed,
			entities.ResponseMessage{"Error: invalid branch"},
			http.MethodPatch,
			"Student should be updated and status code should be 200",
		},
	}

	for i := range testcases {
		var req *http.Request
		switch testcases[i].method {
		case http.MethodPost, http.MethodPut:
			reqBody, _ := json.Marshal(testcases[i].body)
			req = httptest.NewRequest(testcases[i].method, URL, bytes.NewBuffer(reqBody))
		case http.MethodDelete:
			req = httptest.NewRequest(testcases[i].method, URL+"?id="+testcases[i].body.(string), nil)
		case http.MethodGet:
			req = httptest.NewRequest(testcases[i].method, URL+"?id="+testcases[i].body.(bodyStruct).id+
				"&name="+testcases[i].body.(bodyStruct).name+"&branch="+testcases[i].body.(bodyStruct).branch+
				"&includeCompany="+testcases[i].body.(bodyStruct).includeCompany, nil)
		default:
			req = httptest.NewRequest(testcases[i].method, URL, nil)
		}
		w := httptest.NewRecorder()
		handler := New(mockStudentService{})

		handler.Handler(w, req)

		if w.Code != testcases[i].expecStatus {
			t.Errorf(" Test: %v\n Expected Code: %v\n Actual Code: %v\n Description: %v", i+1,
				testcases[i].expecStatus, w.Code, testcases[i].description)
		}
	}
}

// TestHandler_Get test function to test Student Get handler function
func TestHandler_Get(t *testing.T) {
	type bodyStruct struct {
		name           string
		branch         string
		includeCompany string
	}

	testcases := []struct {
		body          bodyStruct
		expecStatus   int
		expecResponse []entities.Student
		description   string
	}{
		{
			bodyStruct{
				"Student",
				"CSE",
				"false",
			},
			http.StatusOK,
			[]entities.Student{
				{"1", "Student", "12/12/2000", "CSE",
					"9876543210", entities.Company{ID: "1"}, "Pending"},
			},
			"Student with that name and branch is present so a student should be returned and status code should be 200",
		},
		{
			bodyStruct{
				"Student",
				"CSE",
				"true",
			},
			http.StatusOK,
			[]entities.Student{
				{"1", "Student", "12/12/2000", "CSE",
					"9876543210", entities.Company{"1", "Test Company", "MASS"}, "Pending"},
			},
			"Student with that name and branch is present includeCompany flag is true so a student " +
				"with company detail should be returned and status code should be 200",
		},
		{
			bodyStruct{
				"Student",
				"CSE2",
				"false",
			},
			http.StatusNotFound,
			[]entities.Student{},
			"Student with that name and branch branch is not present so empty json object should be returned " +
				"with status code 404",
		},
		{
			bodyStruct{
				"Student5",
				"CSE",
				"false",
			},
			http.StatusNotFound,
			[]entities.Student{},
			"Student with that name and branch branch is not present so empty json object should be returned" +
				" with status code 404",
		},
		{
			bodyStruct{
				"Student5",
				"CSE",
				"random",
			},
			http.StatusBadRequest,
			[]entities.Student{},
			"includeCompany can only be either true or false so status code is 400",
		},
	}

	for i := range testcases {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("%v?name=%v&branch=%v&includeCompany=%v", URL,
			testcases[i].body.name, testcases[i].body.branch, testcases[i].body.includeCompany), nil)
		w := httptest.NewRecorder()
		handler := New(mockStudentService{})

		handler.Get(w, req)

		if w.Code != testcases[i].expecStatus {
			t.Errorf(" Test: %v\n Expected Code: %v\n Actual Code: %v\n Description: %v", i+1,
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
				"9876543210", entities.Company{ID: "1"}, "Pending"},
			"Student with that ID is present so a student should be returned and status code should be 200",
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

		handler.GetByID(w, req)

		if w.Code != testcases[i].expecStatus {
			t.Errorf(" Test: %v\n Expected Code: %v\n Actual Code: %v\n Description: %v", i+1,
				testcases[i].expecStatus, w.Code, testcases[i].description)
		}
	}
}

// TestHandler_Create test function to test Student Create handler
func TestHandler_Create(t *testing.T) {
	testcases := []struct {
		body        interface{}
		expStatus   int
		description string
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
			http.StatusCreated,
			"Student should be added and status code should be 201",
		},
		{
			entities.Student{},
			http.StatusBadRequest,
			"No details of student passed so status code should be 400",
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
			http.StatusBadRequest,
			"Student should not be created as branch is not valid and status code should be 400",
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
			http.StatusBadRequest,
			"Student should not be created as phone is not valid and status code should be 400",
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
			http.StatusBadRequest,
			"Student should not be created as status is not valid and status code should be 400",
		},
		{
			struct {
				RandomParam string `json:"random_param"`
			}{
				"random value",
			},
			http.StatusBadRequest,
			"Malformed request body and status code should be 400",
		},
	}

	for i := range testcases {
		reqBody, _ := json.Marshal(testcases[i].body)
		req := httptest.NewRequest(http.MethodPost, URL, bytes.NewReader(reqBody))
		w := httptest.NewRecorder()
		handler := New(mockStudentService{})

		handler.Create(w, req)

		if w.Code != testcases[i].expStatus {
			t.Errorf(" Test: %v\n Expected Code: %v\n Actual Code: %v\n Description: %v", i+1,
				testcases[i].expStatus, w.Code, testcases[i].description)
		}
	}
}

// TestHandler_Update test function to test Student Update handler
func TestHandler_Update(t *testing.T) {
	testcases := []struct {
		body        entities.Student
		expecStatus int
		description string
	}{
		{
			entities.Student{
				"1",
				"Test Student",
				"12/12/2000",
				"ECE",
				"9876543210",
				entities.Company{ID: "1"},
				"CORE",
			},
			http.StatusOK,
			"Student should be updated and status code should be 200",
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
			http.StatusBadRequest,
			"Student should not be update as branch is not valid and status code should be 400",
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
			http.StatusBadRequest,
			"Student should not be update as phone no. is not valid and status code should be 400",
		},
		{
			entities.Student{
				"1",
				"Test Student",
				"12/12/2000",
				"ECE",
				"987654321013311",
				entities.Company{ID: "1"},
				"CORE",
			},
			http.StatusBadRequest,
			"Student should not be update as status is not valid and status code should be 400",
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
			http.StatusNotFound,
			"Student should not be update as no student with this id and status code should be 404",
		},
		{
			entities.Student{ID: "1"},
			http.StatusBadRequest,
			"Student should not be update as no student with this id and status code should be 404",
		},
	}

	for i := range testcases {
		reqBody, _ := json.Marshal(testcases[i].body)
		req := httptest.NewRequest(http.MethodPut, URL+"?id="+testcases[i].body.ID, bytes.NewReader(reqBody))
		w := httptest.NewRecorder()
		handler := New(mockStudentService{})

		handler.Update(w, req)

		if w.Code != testcases[i].expecStatus {
			t.Errorf(" Test: %v\n Expected Code: %v\n Actual Code: %v\n Description: %v", i+1,
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
		{
			"",
			http.StatusBadRequest,
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
			t.Errorf(" Test: %v\n Expected Code: %v\n Actual Code: %v\n Description: %v", i+1,
				testcases[i].expecStatus, w.Code, testcases[i].description)
		}
	}
}

type mockStudentService struct{}

// Get mock services for Get for Student
func (m mockStudentService) Get(ctx context.Context, name string, branch string,
	includeCompany bool) ([]entities.Student, error) {
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
	return []entities.Student{}, errors.EntityNotFound{"Student not found"}
}

// GetByID mock services for GetByID for Student
func (m mockStudentService) GetByID(ctx context.Context, id string) (entities.Student, error) {
	if id != "1" {
		return entities.Student{}, errors.EntityNotFound{"Student"}
	}
	return entities.Student{"1", "Test Student 1", "12/12/2000", "CSE", "9876543210",
		entities.Company{}, "PENDING"}, nil
}

// Create mock service for Create of Student
func (m mockStudentService) Create(ctx context.Context, student entities.Student) (entities.Student, error) {
	if len(student.Phone) < 10 || len(student.Phone) > 12 {
		return entities.Student{}, errors.InvalidParams{"Phone must be of 10-12 digits"}
	}

	if !(student.Branch == "CSE" || student.Branch == "ISE" || student.Branch == "MECH" || student.Branch == "CIVIL" ||
		student.Branch == "ECE" || student.Branch == "EEE") {
		return entities.Student{}, errors.InvalidParams{"invalid branch"}
	}

	if !(student.Status == "PENDING" || student.Status == "ACCEPTED" || student.Status == "REJECTED") {
		return entities.Student{}, errors.InvalidParams{"invalid status"}
	}

	return entities.Student{
		Name:    "Test Student",
		DOB:     "12/12/2000",
		Phone:   "9876543210",
		Branch:  "CSE",
		Company: entities.Company{ID: "1"},
		Status:  "PENDING",
	}, nil
}

// Update mock service for Update of Student
func (m mockStudentService) Update(ctx context.Context, student entities.Student) (entities.Student, error) {
	if student.ID != "1" {
		return entities.Student{}, errors.EntityNotFound{"Student"}
	}

	if len(student.Phone) < 10 || len(student.Phone) > 12 {
		return entities.Student{}, errors.InvalidParams{"invalid phone no."}
	}

	if !(student.Branch == "CSE" || student.Branch == "ISE" || student.Branch == "MECH" || student.Branch == "CIVIL" ||
		student.Branch == "ECE" || student.Branch == "EEE") {
		return entities.Student{}, errors.InvalidParams{"invalid branch"}
	}

	if student.Status == "PENDING" || student.Status == "ACCEPTED" || student.Status == "REJECTED" {
		return entities.Student{}, errors.InvalidParams{"invalid status"}
	}

	return entities.Student{}, nil
}

// Delete mock service for Delete of Student
func (m mockStudentService) Delete(ctx context.Context, id string) error {
	if id != "1" {
		return errors.EntityNotFound{"Student"}
	}
	return nil
}
