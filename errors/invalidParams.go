package errors

import "fmt"

type InvalidParams struct {
	Message string `json:"message"`
}

func (i InvalidParams) Error() string {
	return fmt.Sprintf("Error: %v", i.Message)
}
