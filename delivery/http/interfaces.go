package http

import (
	"net/http"
)

type CompanyInterface interface {
	GetHandler(w http.ResponseWriter, req *http.Request)
	CreateHandler(w http.ResponseWriter, req *http.Request)
	UpdateHandler(w http.ResponseWriter, req *http.Request)
	DeleteHandler(w http.ResponseWriter, req *http.Request) (int, error)
}

type StudentInterface interface {
	GetHandler(w http.ResponseWriter, req *http.Request)
	GetByIdHandler(w http.ResponseWriter, req *http.Request)
	SearchHandler(w http.ResponseWriter, req *http.Request)
	UpdateHandler(w http.ResponseWriter, req *http.Request)
	DeleteHandler(w http.ResponseWriter, req *http.Request)
}
