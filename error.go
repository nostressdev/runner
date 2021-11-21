package runner

import "errors"

var NoDeadline = errors.New("deadline must be set")
var InitializationDeadline = errors.New("initialization deadline")
var ShutdownDeadline = errors.New("graceful shutdown deadline")

type MultiError struct {
	errors []error
}

func NewMultiError(errors ...error) error {
	result := make([]error, 0)
	for _, err := range errors {
		if err != nil {
			result = append(result, err)
		}
	}
	if len(result) == 0 {
		return nil
	}
	return &MultiError{
		errors: result,
	}
}

func (multi *MultiError) Error() string {
	result := ""
	for _, err := range multi.errors {
		result += err.Error()
	}
	return result
}
