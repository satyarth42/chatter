package dto

import "fmt"

type CommonError struct {
	Err        error
	StatusCode int
}

func (e *CommonError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s, statusCode: %d", e.Err.Error(), e.StatusCode)
	}
	return fmt.Sprintf("statusCode: %d", e.StatusCode)
}
