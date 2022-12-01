package errors

import "fmt"

type MissisngParam struct {
	Params []string `json:"param"`
}

func (m MissisngParam) Error() string {
	return fmt.Sprintf("Required params are missing. Params: %v", m.Params)
}
