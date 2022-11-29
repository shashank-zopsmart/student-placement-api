package store

import "student-placement-api/entities"

type Company interface {
	GetByID(id string) (entities.Company, error)
	Create(company entities.Company) (entities.Company, error)
	Update(company entities.Company) (entities.Company, error)
	Delete(id string) (entities.Company, error)
}

type Student interface {
	Get(name string, branch string, includeCompany bool) ([]entities.Student, error)
	GetById(id string) (entities.Student, error)
	Create(student entities.Student) (entities.Student, error)
	Update(student entities.Student) (entities.Student, error)
	Delete(id string) (entities.Student, error)
}
