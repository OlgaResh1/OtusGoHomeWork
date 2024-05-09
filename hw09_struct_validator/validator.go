package hw09structvalidator

import (
	"reflect"
)

func Validate(v interface{}) error {
	validator := StructValidator{
		validateType:  reflect.TypeOf(v),
		validateValue: reflect.ValueOf(v),
	}

	if err := validator.ValidateStruct(); err != nil {
		return err
	}
	if len(validator.validateErrs) > 0 {
		return validator.validateErrs
	}
	return nil
}
