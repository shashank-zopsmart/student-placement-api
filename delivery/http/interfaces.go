package http

import (
	"net/http"
)

type Company interface {
	Handler(w http.ResponseWriter, req *http.Request)
	GetByID(w http.ResponseWriter, req *http.Request)
	Create(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
	Delete(w http.ResponseWriter, req *http.Request)
}

type Student interface {
	Handler(w http.ResponseWriter, req *http.Request)
	Get(w http.ResponseWriter, req *http.Request)
	GetByIdHandler(w http.ResponseWriter, req *http.Request)
	Create(w http.ResponseWriter, req *http.Request)
	Update(w http.ResponseWriter, req *http.Request)
	Delete(w http.ResponseWriter, req *http.Request)
}
