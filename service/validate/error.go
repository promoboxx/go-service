package validate

import "fmt"

type ValidationError interface {
	InvalidField() string
	InnerError() error
	Error() string
}

func NewValidationError(invalidField string, innerError error) ValidationError {
	return &validationError{invalidField: invalidField, innerError: innerError}
}

type validationError struct {
	invalidField string
	innerError   error
}

func (ve *validationError) InvalidField() string {
	return ve.invalidField
}

func (ve *validationError) InnerError() error {
	return ve.innerError
}

func (ve *validationError) Error() string {
	return fmt.Sprintf("Invalid Field: %s- %v", ve.invalidField, ve.innerError)
}
