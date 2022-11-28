package service

import "student-placement-api/entities"

type Company interface {
	GetByID(id string) (entities.Company, error)
	Create(company entities.Company) error
	Update(company entities.Company) error
	Delete(id string) error
}

type Student interface {
	Get(name string, branch string, includeCompany bool) ([]entities.Student, error)
	GetById(id string) (entities.Student, error)
	Create(student entities.Student) error
	Update(student entities.Student) error
	Delete(id string) error
}
