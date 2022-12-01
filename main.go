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

	http.HandleFunc("/company", companyHandler.Handler)
	http.HandleFunc("/student", studentHandler.Handler)

	log.Fatal(http.ListenAndServe(":8080", nil))
}
