package errors

import "fmt"

type HttpErrors struct {
	ErrorMessage string `json:"error"`
}

func (h HttpErrors) Error() string {
	return fmt.Sprintf("%v", h.ErrorMessage)
}