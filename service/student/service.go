package student

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
	"student-placement-api/entities"
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
func (service service) Create(student entities.Student) (entities.Student, error) {
	if len(student.Name) < 3 {
		return entities.Student{}, errors.New("invalid name, it must of minimum 3 characters")
	}

	if ageValidate, err := validateAge(student.DOB, 22); err != nil {
		return entities.Student{}, errors.New("invalid dob, use dd/mm/yyyy")
	} else if ageValidate == false {
		return entities.Student{}, errors.New("student doesn't meet minimum age requirement")
	}

	if len(student.Phone) < 10 || len(student.Phone) > 12 {
		return entities.Student{}, errors.New("invalid phone, must be of 10-12 digits")
	}

	if validatePhone(student.Phone) == false {
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
		return entities.Student{}, errors.New(fmt.Sprintf("Company not found %v", err))
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

// Get service to get all student or search student by name and branch
func (service service) Get(name string, branch string, includeCompany bool) ([]entities.Student, error) {
	return service.store.Get(name, branch, includeCompany)
}

// GetByID service to get a student by ID
func (service service) GetByID(id string) (entities.Student, error) {
	return service.store.GetById(id)
}

// Update service to update a particular student
func (service service) Update(student entities.Student) (entities.Student, error) {
	_, err := service.store.GetById(student.ID)
	if err != nil {
		return entities.Student{}, err
	}

	if len(student.Name) < 3 {
		return entities.Student{}, errors.New("invalid name, it must of minimum 3 characters")
	}

	if ageValidate, err := validateAge(student.DOB, 22); err != nil {
		return entities.Student{}, errors.New("invalid dob, use dd/mm/yyyy")
	} else if ageValidate == false {
		return entities.Student{}, errors.New("student doesn't meet minimum age requirement")
	}

	if len(student.Phone) < 10 || len(student.Phone) > 12 {
		return entities.Student{}, errors.New("invalid phone, must be of 10-12 digits")
	}

	if validatePhone(student.Phone) == false {
		return entities.Student{}, errors.New("invalid phone")
	}

	if !(student.Status == "PENDING" || student.Status == "ACCEPTED" || student.Status == "REJECTED") {
		return entities.Student{}, errors.New("invalid status")
	}

	if !(student.Branch == "CSE" || student.Branch == "ISE" || student.Branch == "MECH" || student.Branch == "CIVIL" ||
		student.Branch == "ECE" || student.Branch == "EEE") {
		return entities.Student{}, errors.New("invalid branch")
	}

	company, err := service.store.GetCompany(student.Company.ID)

	if err != nil {
		return entities.Student{}, err
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

	return service.store.Update(student)
}

// Delete service to delete a particular student
func (service service) Delete(id string) error {
	_, err := service.store.GetById(id)
	if err != nil {
		return err
	}
	return service.store.Delete(id)
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
