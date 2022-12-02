package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	companyHandler "student-placement-api/delivery/http/company"
	studentHandler "student-placement-api/delivery/http/student"
	companyService "student-placement-api/service/company"
	studentService "student-placement-api/service/student"
	_ "student-placement-api/store"
	"student-placement-api/store/company"
	"student-placement-api/store/student"
)

func main() {
	db, err := sql.Open("mysql", "zopsmart:Zopsmart123@tcp(127.0.0.1:3306)/placement")
	if err != nil {
		log.Fatal(err)
	}

	companyStore := company.New(db)
	studentStore := student.New(db)

	companyService := companyService.New(companyStore)
	studentService := studentService.New(studentStore)

	companyHandler := companyHandler.New(companyService)
	studentHandler := studentHandler.New(studentService)

	http.Handle("/company", middleware(http.HandlerFunc(companyHandler.Handler)))
	http.Handle("/student", middleware(http.HandlerFunc(studentHandler.Handler)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func middleware(originalHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		if !(req.Header.Get("x-api-key") == "a601e44e306e430f8dde987f65844f05" ||
			req.Header.Get("x-api-key") == "84dcb7c09b4a4af8a67f4577ffe9b255") {
			return
		}

		if req.Header.Get("Content-Type") != "application/json" {
			w.WriteHeader(http.StatusUnsupportedMediaType)
			w.Write([]byte("Header Content-Type incorrect"))
			return
		}
		originalHandler.ServeHTTP(w, req)
	})
}
