package errors

import "fmt"

type EntityNotFound struct {
	Entity string `json:"entity"`
}

func (e EntityNotFound) Error() string {
	return fmt.Sprintf("%v not found", e.Entity)
}

type ConnDone struct{}

func (e ConnDone) Error() string {
	return fmt.Sprintf("DB connection closed")
}
