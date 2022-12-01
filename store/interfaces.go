package store

import (
	"context"
	"student-placement-api/entities"
)

type Company interface {
	GetByID(ctx context.Context, id string) (entities.Company, error)
	Create(ctx context.Context, company entities.Company) (entities.Company, error)
	Update(ctx context.Context, company entities.Company) (entities.Company, error)
	Delete(ctx context.Context, id string) error
}

type Student interface {
	Get(ctx context.Context, name string, branch string, includeCompany bool) ([]entities.Student, error)
	GetById(ctx context.Context, id string) (entities.Student, error)
	Create(ctx context.Context, student entities.Student) (entities.Student, error)
	Update(ctx context.Context, student entities.Student) (entities.Student, error)
	Delete(ctx context.Context, id string) error
	GetCompany(ctx context.Context, id string) (entities.Company, error)
}
