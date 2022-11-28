package store

import "student-placement-api/entities"

type Company interface {
	Get(id string) (entities.Company, error)
	Create(id string, name string) (int, error)
	Update(company entities.Company) (int, error)
	Delete(id string) (int, error)
}

type Student interface {
	Get() ([]entities.Student, error)
	GetById(id string, includeCompany bool) (entities.Student, error)
	Search(id string, branch string, includeCompany bool) (entities.Student, error)
	Update(student entities.Student) (int, error)
	Delete(id string) (int, error)
}
