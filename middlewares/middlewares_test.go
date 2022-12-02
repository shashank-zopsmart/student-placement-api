package middlewares

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

// Test_Middleware function to test middleware
func Test_Middleware(t *testing.T) {
	testcases := []struct {
		x_api_key    string
		content_type string
		method       string
		expCode      int
		desc         string
	}{
		{"a601e44e306e430f8dde987f65844f05", "application/json", http.MethodGet,
			http.StatusOK, "x-api-key and content-type both are valid so status code 200"},
		{"84dcb7c09b4a4af8a67f4577ffe9b255", "application/json", http.MethodGet,
			http.StatusOK, "x-api-key and content-type both are valid so status code 200"},
		{"a601e44e306e430f8dde987f65844f05", "application/json", http.MethodPost,
			http.StatusOK, "x-api-key and content-type both are valid so status code 200"},
		{"randomkey", "application/json", http.MethodGet,
			http.StatusForbidden, "x-api-key is invalid so status code 403"},
		{"a601e44e306e430f8dde987f65844f05", "text/plain", http.MethodPost,
			http.StatusUnsupportedMediaType, "Content-Type is invalid so status code 415"},
		{content_type: "text/plain", method: http.MethodGet, expCode: http.StatusForbidden,
			desc: "x-api-key is not present so status code is 403"},
		{x_api_key: "a601e44e306e430f8dde987f65844f05", method: http.MethodGet, expCode: http.StatusUnsupportedMediaType,
			desc: "Content-Type is not present so status code is 415"},
	}

	for i, _ := range testcases {
		req := httptest.NewRequest(testcases[i].method, "http://localhost:8080/", nil)
		w := httptest.NewRecorder()

		req.Header.Set("x-api-key", testcases[i].x_api_key)
		req.Header.Set("Content-Type", testcases[i].content_type)

		handle := Middleware(http.HandlerFunc(mockHandler))
		handle.ServeHTTP(w, req)

		if w.Code != testcases[i].expCode {
			t.Errorf("Test: %v\t Expected Code: %v\t Actual Code: %v\t Description: %v", i+1,
				testcases[i].expCode, w.Code, testcases[i].desc)
		}
	}
}

func mockHandler(w http.ResponseWriter, req *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(fmt.Sprintf("Called with: %v", req.Method)))
}
