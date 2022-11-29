package student

import (
	"errors"
	"student-placement-api/entities"
	"student-placement-api/store"
)

type Service struct {
	store store.Student
}

// New factory function to return service object and do dependency injection
func New(store store.Student) Service {
	return Service{store: store}
}

// Get service to get all student or search student by name and branch
func (service Service) Get(name string, branch string, includeCompany bool) ([]entities.Student, error) {
	return service.store.Get(name, branch, includeCompany)
}

// GetByID service to get a student by ID
func (service Service) GetByID(id string) (entities.Student, error) {
	return service.store.GetById(id)
}

// Create to create a new student
func (service Service) Create(student entities.Student) (entities.Student, error) {
	if len(student.Name) < 3 {
		return entities.Student{}, errors.New("invalid name")
	}

	if len(student.Phone) < 10 || len(student.Phone) > 12 {
		return entities.Student{}, errors.New("invalid phone")
	}

	if !(student.Status == "PENDING" || student.Status == "ACCEPTED" || student.Status == "REJECTED") {
		return entities.Student{}, errors.New("invalid status")
	}

	if !(student.Branch == "CSE" || student.Branch == "ISE" || student.Branch == "MECH" || student.Branch == "CIVIL" ||
		student.Branch == "ECE" || student.Branch == "EEE") {
		return entities.Student{}, errors.New("invalid branch")
	}

	var company, err = service.store.GetCompany(student.Company.ID)

	if err != nil {
		return entities.Student{}, errors.New("invalid company")
	}

	if company.Category == "DREAM IT" && !(student.Branch == "CSE" || student.Branch == "ISE") {
		return entities.Student{}, errors.New("branch not allowed in this company")
	}

	if company.Category == "OPEN DREAM" && !(student.Branch == "CSE" || student.Branch == "ISE" ||
		student.Branch == "ECE" || student.Branch == "EEE") {
		return entities.Student{}, errors.New("branch not allowed in this company")
	}

	if company.Category == "CORE" && !(student.Branch == "CIVIL" || student.Branch == "MECH") {
		return entities.Student{}, errors.New("branch not allowed in this company")
	}

	return service.store.Create(student)
}

// Update service to update a particular student
func (service Service) Update(student entities.Student) (entities.Student, error) {
	if len(student.Name) < 3 {
		return entities.Student{}, errors.New("invalid name")
	}
	if len(student.Phone) < 10 || len(student.Phone) > 12 {
		return entities.Student{}, errors.New("invalid phone")
	}
	if !(student.Status == "PENDING" || student.Status == "ACCEPTED" || student.Status == "REJECTED") {
		return entities.Student{}, errors.New("invalid status")
	}

	if !(student.Branch == "CSE" || student.Branch == "ISE" || student.Branch == "MECH" || student.Branch == "CIVIL" ||
		student.Branch == "ECE" || student.Branch == "EEE") {
		return entities.Student{}, errors.New("invalid branch")
	}

	var company, err = service.store.GetCompany(student.Company.ID)

	if err != nil {
		return entities.Student{}, errors.New("invalid company")
	}

	if company.Category == "DREAM IT" && !(student.Branch == "CSE" || student.Branch == "ISE") {
		return entities.Student{}, errors.New("branch not allowed in this company")
	}

	if company.Category == "OPEN DREAM" && !(student.Branch == "CSE" || student.Branch == "ISE" ||
		student.Branch == "ECE" || student.Branch == "EEE") {
		return entities.Student{}, errors.New("branch not allowed in this company")
	}

	if company.Category == "CORE" && !(student.Branch == "CIVIL" || student.Branch == "MECH") {
		return entities.Student{}, errors.New("branch not allowed in this company")
	}

	return service.store.Create(student)
}

// Delete service to delete a particular student
func (service Service) Delete(id string) (entities.Student, error) {
	return service.store.Delete(id)
}
