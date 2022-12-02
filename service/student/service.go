package student

import (
	"context"
	"strconv"
	"strings"
	"student-placement-api/entities"
	"student-placement-api/errors"
	"student-placement-api/store"
	"time"
	"unicode"
)

type service struct {
	store store.Student
}

// New factory function to return service object and do dependency injection
func New(store store.Student) service {
	return service{store: store}
}

// Create to create a new student
func (service service) Create(ctx context.Context, student entities.Student) (entities.Student, error) {
	if len(student.Name) < 3 {
		return entities.Student{}, errors.InvalidParams{Message: "Name should be of at least 2 characters"}
	}

	if ageValidate, err := validateAge(student.DOB, 22); err != nil {
		return entities.Student{}, errors.InvalidParams{Message: "Invalid DOB, Use dd/mm/yyyy"}
	} else if ageValidate == false {
		return entities.Student{}, errors.InvalidParams{Message: "Student doesn't meet minimum age requirement"}
	}

	if len(student.Phone) < 10 || len(student.Phone) > 12 {
		return entities.Student{}, errors.InvalidParams{Message: "Phone must be of 10-12 digits"}
	}

	if validatePhone(student.Phone) == false {
		return entities.Student{}, errors.InvalidParams{Message: "Invalid Phone"}
	}

	if !(student.Status == "PENDING" || student.Status == "ACCEPTED" || student.Status == "REJECTED") {
		return entities.Student{}, errors.InvalidParams{Message: "Invalid Status"}
	}

	if !(student.Branch == "CSE" || student.Branch == "ISE" || student.Branch == "MECH" || student.Branch == "CIVIL" ||
		student.Branch == "ECE" || student.Branch == "EEE") {
		return entities.Student{}, errors.InvalidParams{Message: "Invalid Branch"}
	}

	var company, err = service.store.GetCompany(ctx, student.Company.ID)
	if err != nil {
		return entities.Student{}, errors.EntityNotFound{Entity: "Company"}
	}

	if company.Category == "DREAM IT" && !(student.Branch == "CSE" || student.Branch == "ISE") {
		return entities.Student{}, errors.InvalidParams{Message: "Branch not allowed in this company"}
	}

	if company.Category == "OPEN DREAM" && !(student.Branch == "CSE" || student.Branch == "ISE" ||
		student.Branch == "ECE" || student.Branch == "EEE") {
		return entities.Student{}, errors.InvalidParams{Message: "Branch not allowed in this company"}
	}

	if company.Category == "CORE" && !(student.Branch == "CIVIL" || student.Branch == "MECH") {
		return entities.Student{}, errors.InvalidParams{Message: "Branch not allowed in this company"}
	}

	return service.store.Create(ctx, student)
}

// Get service to get all student or search student by name and branch
func (service service) Get(ctx context.Context, name string, branch string, includeCompany bool) ([]entities.Student, error) {
	return service.store.Get(ctx, name, branch, includeCompany)
}

// GetByID service to get a student by ID
func (service service) GetByID(ctx context.Context, id string) (entities.Student, error) {
	return service.store.GetById(ctx, id)
}

// Update service to update a particular student
func (service service) Update(ctx context.Context, student entities.Student) (entities.Student, error) {
	_, err := service.store.GetById(ctx, student.ID)
	if err != nil {
		return entities.Student{}, errors.EntityNotFound{Entity: "Student"}
	}

	if len(student.Name) < 3 {
		return entities.Student{}, errors.InvalidParams{Message: "Name should be of at least 2 characters"}
	}

	if ageValidate, err := validateAge(student.DOB, 22); err != nil {
		return entities.Student{}, errors.InvalidParams{Message: "Invalid DOB, Use dd/mm/yyyy"}
	} else if ageValidate == false {
		return entities.Student{}, errors.InvalidParams{Message: "Student doesn't meet minimum age requirement"}
	}

	if len(student.Phone) < 10 || len(student.Phone) > 12 {
		return entities.Student{}, errors.InvalidParams{Message: "Phone must be of 10-12 digits"}
	}

	if validatePhone(student.Phone) == false {
		return entities.Student{}, errors.InvalidParams{Message: "Invalid Phone"}
	}

	if !(student.Status == "PENDING" || student.Status == "ACCEPTED" || student.Status == "REJECTED") {
		return entities.Student{}, errors.InvalidParams{Message: "Invalid Status"}
	}

	if !(student.Branch == "CSE" || student.Branch == "ISE" || student.Branch == "MECH" || student.Branch == "CIVIL" ||
		student.Branch == "ECE" || student.Branch == "EEE") {
		return entities.Student{}, errors.InvalidParams{Message: "Invalid Branch"}
	}

	company, err := service.store.GetCompany(ctx, student.Company.ID)
	if err != nil {
		return entities.Student{}, errors.EntityNotFound{Entity: "Company"}
	}

	if company.Category == "DREAM IT" && !(student.Branch == "CSE" || student.Branch == "ISE") {
		return entities.Student{}, errors.InvalidParams{Message: "Branch not allowed in this company"}
	}

	if company.Category == "OPEN DREAM" && !(student.Branch == "CSE" || student.Branch == "ISE" ||
		student.Branch == "ECE" || student.Branch == "EEE") {
		return entities.Student{}, errors.InvalidParams{Message: "Branch not allowed in this company"}
	}

	if company.Category == "CORE" && !(student.Branch == "CIVIL" || student.Branch == "MECH") {
		return entities.Student{}, errors.InvalidParams{Message: "Branch not allowed in this company"}
	}

	return service.store.Update(ctx, student)
}

// Delete service to delete a particular student
func (service service) Delete(ctx context.Context, id string) error {
	_, err := service.store.GetById(ctx, id)
	if err != nil {
		return errors.EntityNotFound{Entity: "Student"}
	}
	return service.store.Delete(ctx, id)
}

func validateAge(dob string, minAge int) (bool, error) {
	year := time.Now().Year()
	dob = strings.TrimSpace(dob)

	yob, err := strconv.Atoi(string(dob[len(dob)-4:]))
	if err != nil {
		return false, err
	}
	age := year - yob
	if age >= minAge {
		return true, nil
	}
	return false, nil
}

func validatePhone(phone string) bool {
	for _, char := range phone {
		if !unicode.IsNumber(char) {
			return false
		}
	}
	return true
}
