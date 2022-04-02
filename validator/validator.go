package validator

import (
	"errors"
	"github.com/go-playground/validator/v10"
	"github.com/oscto/ky3k"
	"strings"
)

type Validator struct {
	validator *validator.Validate
}

func NewValidator() *Validator {

	return &Validator{
		validator: validator.New(),
	}
}

func (v *Validator) Bind(str string, data interface{}) error {

	if err := ky3k.StringToJson(str, &data); err != nil {
		return err
	}

	var errStr []string
	if err := v.validator.Struct(data); err != nil {
		if errs, ok := err.(validator.ValidationErrors); ok {
			for _, fieldError := range errs {
				errStr = append(errStr, fieldError.Error())
			}
		} else {
			errStr = append(errStr, err.Error())
		}
	}
	if len(errStr) > 0 {
		return errors.New(strings.Join(errStr, ";"))
	}

	return nil
}
