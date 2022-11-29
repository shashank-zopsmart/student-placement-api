package company

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"student-placement-api/entities"
	"testing"
)

const URL = "http://localhost:8080/company"

// TestHandler_Handler test function to test Company main handler function
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
			entities.Company{"1", "Test Company", "MASS"},
			http.MethodGet,
			"Company with that ID is present so a company should be returned and status code should be 200",
		},
		{
			entities.Company{
				Name:     "Test Company",
				Category: "MASS",
			},
			http.StatusCreated,
			nil,
			http.MethodPost,
			"Company should be added and status code should be 201",
		},
		{
			entities.Company{
				"1",
				"Test Company",
				"MASS",
			},
			http.StatusOK,
			entities.ResponseMessage{"Company Updated"},
			http.MethodPut,
			"Company should be updated and status code should be 200",
		},
		{
			"1",
			http.StatusOK,
			entities.ResponseMessage{"Company deleted"},
			http.MethodDelete,
			"Company with that ID should be deleted and status code should be 200",
		},
	}

	for i := range testcases {
		var req *http.Request
		switch testcases[i].method {
		case http.MethodPost:
			reqBody, _ := json.Marshal(testcases[i].body)
			req = httptest.NewRequest(testcases[i].method, URL, bytes.NewBuffer(reqBody))
		case http.MethodPut:
			reqBody, _ := json.Marshal(testcases[i].body)
			req = httptest.NewRequest(testcases[i].method, URL+"?id="+testcases[i].body.(entities.Company).ID,
				bytes.NewBuffer(reqBody))
		default:
			req = httptest.NewRequest(testcases[i].method, URL+"?id="+testcases[i].body.(string), nil)
		}
		w := httptest.NewRecorder()
		handler := New(mockCompanyService{})

		handler.Handler(w, req)

		if w.Code != testcases[i].expecStatus {
			t.Errorf("Test: %v\t Expected Code: %v\t Actual Code: %v\t Description: %v", i+1,
				testcases[i].expecStatus, w.Code, testcases[i].description)
		}
	}
}

// TestHandler_Get test function to test Company Get handler function
func TestHandler_Get(t *testing.T) {
	testcases := []struct {
		id            string
		expecStatus   int
		expecResponse interface{}
		description   string
	}{
		{
			"1",
			http.StatusOK,
			entities.Company{"1", "Test Company", "MASS"},
			"Company with that ID is present so a company should be returned and status code should be 200",
		},
		{
			"2",
			http.StatusNotFound,
			entities.Company{},
			"Company with that ID is not present so empty json object should be returned wit status code 404",
		},
		{
			"",
			http.StatusBadRequest,
			entities.Company{},
			"Company with that ID is not present so empty json object should be returned wit status code 404",
		},
	}

	for i := range testcases {
		req := httptest.NewRequest(http.MethodGet, URL+"?id="+testcases[i].id, nil)

		if testcases[i].id == "" {
			req = httptest.NewRequest(http.MethodGet, URL, nil)
		}

		w := httptest.NewRecorder()
		handler := New(mockCompanyService{})

		handler.Get(w, req)

		if w.Code != testcases[i].expecStatus {
			t.Errorf("Test: %v\t Expected Code: %v\t Actual Code: %v\t Description: %v", i+1,
				testcases[i].expecStatus, w.Code, testcases[i].description)
		}
	}
}

// TestHandler_Create test function to test Company Create handler
func TestHandler_Create(t *testing.T) {
	testcases := []struct {
		body          entities.Company
		expecStatus   int
		expecResponse interface{}
		description   string
	}{
		{
			entities.Company{
				Name:     "Test Company",
				Category: "MASS",
			},
			http.StatusCreated,
			entities.Company{
				ID:       "1",
				Name:     "Test Company",
				Category: "MASS",
			},
			"Company should be added and status code should be 201",
		},
		{
			entities.Company{
				Name: "Test Company 2",
			},
			http.StatusBadRequest,
			entities.ResponseMessage{"Error: Name and Category required"},
			"Company should not be created as both parameters are mandatory is not valid and status code " +
				"should be 400",
		},
	}

	for i := range testcases {
		reqBody, _ := json.Marshal(testcases[i].body)
		req := httptest.NewRequest(http.MethodPost, URL, bytes.NewReader(reqBody))
		w := httptest.NewRecorder()
		handler := New(mockCompanyService{})

		handler.Create(w, req)

		if w.Code != testcases[i].expecStatus {
			t.Errorf("Test: %v\t Expected Code: %v\t Actual Code: %v\t Description: %v", i+1,
				testcases[i].expecStatus, w.Code, testcases[i].description)
		}
	}
}

// TestHandler_Update test function to test Company Update handler
func TestHandler_Update(t *testing.T) {
	testcases := []struct {
		body          entities.Company
		expecStatus   int
		expecResponse entities.ResponseMessage
		description   string
	}{
		{
			entities.Company{
				"1",
				"Test Company",
				"MASS",
			},
			http.StatusOK,
			entities.ResponseMessage{"Company Updated"},
			"Company should be updated and status code should be 200",
		},
		{
			entities.Company{
				ID:   "2",
				Name: "Test Company 2",
			},
			http.StatusBadRequest,
			entities.ResponseMessage{"Error: Name and Category required"},
			"Company should not be update as category is missing and status code should be 400",
		},
		{
			entities.Company{
				Name:     "Test Company 3",
				Category: "MASS",
			},
			http.StatusBadRequest,
			entities.ResponseMessage{"Error: ID required"},
			"Company should not be update as id is missing and status code should be 400",
		},
	}

	for i := range testcases {
		reqBody, _ := json.Marshal(testcases[i].body)
		req := httptest.NewRequest(http.MethodPut, URL+"?id="+testcases[i].body.ID, bytes.NewReader(reqBody))
		w := httptest.NewRecorder()
		handler := New(mockCompanyService{})

		handler.Update(w, req)

		if w.Code != testcases[i].expecStatus {
			t.Errorf("Test: %v\t Expected Code: %v\t Actual Code: %v\t Description: %v", i+1,
				testcases[i].expecStatus, w.Code, testcases[i].description)
		}
	}
}

// TestHandler_Delete test function to test Company Delete handler function
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
			entities.ResponseMessage{"Company deleted"},
			"Company with that ID should be deleted and status code should be 200",
		},
		{
			"2",
			http.StatusNotFound,
			entities.ResponseMessage{"Error: Company not found"},
			"Company with that ID is present so a company should be returned and status code should be 200",
		},
		{
			"",
			http.StatusBadRequest,
			entities.ResponseMessage{"Error: ID required"},
			"Id should be passed in the query paramenter and status code should be 200",
		},
	}

	for i := range testcases {
		req := httptest.NewRequest(http.MethodDelete, URL+"?id="+testcases[i].id, nil)
		w := httptest.NewRecorder()
		handler := New(mockCompanyService{})

		handler.Delete(w, req)

		if w.Code != testcases[i].expecStatus {
			t.Errorf("Test: %v\t Expected Code: %v\t Actual Code: %v\t Description: %v", i+1,
				testcases[i].expecStatus, w.Code, testcases[i].description)
		}
	}
}

type mockCompanyService struct{}

// GetByID mock services for GetByID for Company
func (m mockCompanyService) GetByID(id string) (entities.Company, error) {
	if id != "1" {
		return entities.Company{}, errors.New("company not found")
	}
	return entities.Company{"1", "Test Company", "MASS"}, nil
}

// Create mock service for Create of Company
func (m mockCompanyService) Create(company entities.Company) (entities.Company, error) {
	switch company.Category {
	case "MASS", "DREAM IT", "OPEN DREAM", "CORE":
		return entities.Company{"1", "Test Company", "MASS"}, nil
	default:
		return entities.Company{}, errors.New("invalid category")
	}
}

// Update mock service for Update of Company
func (m mockCompanyService) Update(company entities.Company) (entities.Company, error) {
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

// Delete mock service for Delete of Company
func (m mockCompanyService) Delete(id string) (entities.Company, error) {
	if id != "1" {
		return entities.Company{}, errors.New("company not found")
	}
	return entities.Company{}, nil
}
