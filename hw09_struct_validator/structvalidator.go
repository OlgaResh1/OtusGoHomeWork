package hw09structvalidator

import (
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"
)

type StructValidator struct {
	validateErrs  ValidationErrors
	validateType  reflect.Type
	validateValue reflect.Value
	currentField  reflect.StructField
	currentValue  reflect.Value
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

func (v *StructValidator) validateIn(values string) error {
	allowVals := strings.Split(values, ",")

	switch v.currentField.Type.Kind() { //nolint:exhaustive
	case reflect.String:
		{
			for _, allowVal := range allowVals {
				if v.currentValue.String() == allowVal {
					return nil
				}
			}
			v.validateErrs = append(v.validateErrs, ValidationError{
				Field: v.currentField.Name,
				Err:   fmt.Errorf("value must be in list"),
			})
		}
	case reflect.Int:
		{
			for _, allowVal := range allowVals {
				allowInt, err := strconv.Atoi(allowVal)
				if err != nil {
					return UsingError{
						Type: ErrFormatTag, Where: v.currentField.Name,
						Err: fmt.Errorf("error format validate value: %w", err),
					}
				}
				if int(v.currentValue.Int()) == allowInt {
					return nil
				}
			}
			v.validateErrs = append(v.validateErrs, ValidationError{
				Field: v.currentField.Name,
				Err:   fmt.Errorf("value must be in list"),
			})
		}
	case reflect.Slice:
		{
			for i := 0; i < v.currentValue.Len(); i++ {
				val := v.currentValue.Index(i)
				switch {
				case val.CanInt():
					if contains(allowVals, strconv.Itoa(int(val.Int()))) {
						v.validateErrs = append(v.validateErrs, ValidationError{
							Field: v.currentField.Name,
							Err:   fmt.Errorf("value must be in list"),
						})
					}

				case val.Kind() == reflect.String:
					if !contains(allowVals, val.String()) {
						v.validateErrs = append(v.validateErrs, ValidationError{
							Field: v.currentField.Name,
							Err:   fmt.Errorf("value must be in list"),
						})
					}
				default:
					return UsingError{
						Type: ErrUnsupported, Where: v.currentField.Name,
						Err: fmt.Errorf("unsupported validate type: %v", v.currentField.Type),
					}
				}
			}
			return nil
		}
	default:
		{
			return UsingError{
				Type: ErrUnsupported, Where: v.currentField.Name,
				Err: fmt.Errorf("unsupported validate type: %v", v.currentField.Type),
			}
		}
	}
	return nil
}

func (v *StructValidator) validateRange(limit string, isMin bool) error {
	expected, err := strconv.Atoi(limit)
	if err != nil {
		return UsingError{
			Type: ErrFormatTag, Where: v.currentField.Name,
			Err: fmt.Errorf("error format validate value: %w", err),
		}
	}

	switch v.currentField.Type.Kind() { //nolint:exhaustive
	case reflect.Int:
		{
			if int(v.currentValue.Int()) < expected && isMin ||
				int(v.currentValue.Int()) > expected && !isMin {
				v.validateErrs = append(v.validateErrs, ValidationError{
					Field: v.currentField.Name,
					Err:   fmt.Errorf("value is out of range"),
				})
			}
		}
	case reflect.Slice:
		{
			for i := 0; i < v.currentValue.Len(); i++ {
				if !v.currentValue.Index(i).CanInt() {
					return UsingError{
						Type: ErrUnsupported, Where: v.currentField.Name,
						Err: fmt.Errorf("unsupported validate type: %v", v.currentField.Type),
					}
				}
				if int(v.currentValue.Index(i).Int()) < expected && isMin ||
					int(v.currentValue.Index(i).Int()) > expected && !isMin {
					v.validateErrs = append(v.validateErrs, ValidationError{
						Field: v.currentField.Name,
						Err:   fmt.Errorf("value is out of range"),
					})
				}
			}
		}
	default:
		{
			return UsingError{
				Type: ErrUnsupported, Where: v.currentField.Name,
				Err: fmt.Errorf("unsupported validate type: %v", v.currentField.Type),
			}
		}
	}

	return nil
}

func (v *StructValidator) validateMin(minStr string) error {
	return v.validateRange(minStr, true)
}

func (v *StructValidator) validateMax(maxStr string) error {
	return v.validateRange(maxStr, false)
}

func (v *StructValidator) validateRegexp(regexpr string) error {
	re, err := regexp.Compile(regexpr)
	if err != nil {
		return UsingError{
			Type: ErrFormatTag, Where: v.currentField.Name,
			Err: fmt.Errorf("error compile regexp"),
		}
	}

	if !re.MatchString(fmt.Sprint(v.currentValue)) {
		v.validateErrs = append(v.validateErrs, ValidationError{
			Field: v.currentField.Name,
			Err:   fmt.Errorf("value must match regexp"),
		})
	}
	return nil
}

func (v *StructValidator) validateLen(lenString string) error {
	expectedLen, err := strconv.Atoi(lenString)
	if err != nil {
		return UsingError{
			Type: ErrUnsupported, Where: v.currentField.Name,
			Err: fmt.Errorf("unsupported validate type: %v", v.currentField.Type),
		}
	}

	switch v.currentField.Type.Kind() { //nolint:exhaustive
	case reflect.String:
		if expectedLen != len(v.currentValue.String()) {
			v.validateErrs = append(v.validateErrs, ValidationError{
				Field: v.currentField.Name,
				Err:   fmt.Errorf("len does not match %d", len(v.currentValue.String())),
			})
		}
	case reflect.Slice:

		for i := 0; i < v.currentValue.Len(); i++ {
			if v.currentValue.Index(i).Kind() != reflect.String {
				return UsingError{
					Type: ErrUnsupported, Where: v.currentField.Name,
					Err: fmt.Errorf("unsupported validate type: %v", v.currentField.Type),
				}
			}
			if expectedLen != len(v.currentValue.Index(i).String()) {
				v.validateErrs = append(v.validateErrs, ValidationError{
					Field: v.currentField.Name,
					Err:   fmt.Errorf("len does not match %d", len(v.currentValue.Index(i).String())),
				})
			}
		}

	default:
		return UsingError{
			Type: ErrUnsupported, Where: v.currentField.Name,
			Err: fmt.Errorf("unsupported validate type: %v", v.currentField.Type),
		}
	}
	return nil
}

func (v *StructValidator) ValidateStruct() error {
	if v.validateType.Kind() != reflect.Struct {
		return UsingError{
			Type: ErrInvalidType, Where: v.validateType.Name(),
			Err: fmt.Errorf("type is not a struct"),
		}
	}

	funcs := map[string]func(string) error{
		"min":    v.validateMin,
		"max":    v.validateMax,
		"in":     v.validateIn,
		"len":    v.validateLen,
		"regexp": v.validateRegexp,
	}

	for i := 0; i < v.validateType.NumField(); i++ {
		v.currentField = v.validateType.Field(i)
		v.currentValue = v.validateValue.Field(i)
		validate, ok := v.currentField.Tag.Lookup("validate")
		if !ok {
			continue
		}

		for _, cond := range strings.Split(validate, "|") {
			tag, valueTag, found := strings.Cut(cond, ":")
			if !found {
				return UsingError{
					Type: ErrInvalidTag, Where: v.currentField.Name,
					Err: fmt.Errorf("validation tag not contains ':' "),
				}
			}

			method, ok := funcs[tag]
			if ok {
				if err := method(valueTag); err != nil {
					return err
				}
			} else {
				return UsingError{
					Type: ErrInvalidTag, Where: v.currentField.Name,
					Err: fmt.Errorf("validate tag not exist"),
				}
			}
		}
	}
	return nil
}
