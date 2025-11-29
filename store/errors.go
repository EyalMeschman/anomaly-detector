package store

import "fmt"

const (
	unknownErrorString = "ModelStore: unknown error"
	ModuleName         = "ModelStore"
)

type ModelStoreErrorType int

const (
	ErrorInvalidModel ModelStoreErrorType = iota
	ErrorDuplicateModel
	ErrorModelNotFound
)

type ModelStoreError struct {
	ErrType        ModelStoreErrorType
	Err            error
	AdditionalInfo []string
}

func (err *ModelStoreError) Error() string {
	switch err.ErrType {
	case ErrorInvalidModel:
		return "invalid model in batch"
	case ErrorDuplicateModel:
		return fmt.Sprintf("model already exists for %s %s", err.AdditionalInfo[0], err.AdditionalInfo[1])
	case ErrorModelNotFound:
		return fmt.Sprintf("model not found for %s %s", err.AdditionalInfo[0], err.AdditionalInfo[1])
	default:
		if err.Err != nil {
			return fmt.Sprintf("[%s]: %s", unknownErrorString, err.Err.Error())
		}

		return unknownErrorString
	}
}

func isValidationError(err error) bool {
	_, ok := err.(*ModelStoreError)
	return ok
}
