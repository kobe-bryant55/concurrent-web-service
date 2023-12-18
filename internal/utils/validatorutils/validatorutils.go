package validatorutils

import (
	"errors"
	"fmt"
	"github.com/MehmetTalhaSeker/concurrent-web-service/internal/utils/errorutils"

	"github.com/go-playground/validator/v10"
)

type Validator interface {
	Validate(a any) error
}

type customValidator struct {
	validator *validator.Validate
}

func (cv *customValidator) Validate(a interface{}) error {
	err := cv.validator.Struct(a)
	if err != nil {
		var validationErrors validator.ValidationErrors
		if !errors.As(err, &validationErrors) {
			// err is not of type validator.ValidationErrors, return original error
			return err
		}

		var errorList []*errorutils.APIError

		for _, err := range validationErrors {
			switch err.Tag() {
			case "required":
				errorList = append(errorList, errorutils.New(errorutils.Required(err.Field()), err))
			case "max":
				errorList = append(errorList, errorutils.New(errorutils.Max(err.Field()), err))
			case "min":
				errorList = append(errorList, errorutils.New(errorutils.Min(err.Field()), err))
			case "len":
				errorList = append(errorList, errorutils.New(errorutils.Len(err.Field()), err))
			default:
				errorList = append(errorList, errorutils.New(fmt.Errorf("%v doesn't satisfy the constraint", err.Field()), nil))
			}
		}

		return errorutils.ValidationError(errorList)
	}

	return nil
}

func NewValidator() Validator {
	v := validator.New()

	return &customValidator{
		validator: v,
	}
}
