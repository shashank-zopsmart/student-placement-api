package service

import "student-placement-api/entities"

type Company interface {
	GetService(id string) (entities.Company, error)
	CreateService(id string, name string) (int, error)
	UpdateService(company entities.Company) (int, error)
	DeleteService(id string) (int, error)
}

type Student interface {
	GetService() ([]entities.Student, error)
	GetByIdService(id string, includeCompany bool) (entities.Student, error)
	SearchService(id string, branch string, includeCompany bool) (entities.Student, error)
	UpdateService(student entities.Student) (int, error)
	DeleteService(id string) (int, error)
}
