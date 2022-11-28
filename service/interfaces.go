package service

import "student-placement-api/entities"

type Company interface {
	GetByID(id string) entities.Company
	Create(company entities.Company) error
	Update(company entities.Company) error
	Delete(id string) error
}

type Student interface {
	Get(name string, branch string, includeCompany bool) []entities.Student
	GetByID(id string) entities.Student
	Create(student entities.Student) error
	Update(student entities.Student) error
	Delete(id string) error
}
