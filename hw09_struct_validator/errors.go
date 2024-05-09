package hw09structvalidator

import (
	"fmt"
	"strings"
)

type ValidationError struct {
	Field string
	Err   error
}

type ValidationErrors []ValidationError

func (v ValidationError) Error() string {
	return fmt.Sprintf("%s: %s", v.Field, v.Err.Error())
}

func (v ValidationErrors) Error() string {
	var s strings.Builder
	for _, e := range v {
		s.WriteString(e.Field + ": " + e.Err.Error() + "\n")
	}
	return s.String()
}

type UsingErrorType int

const (
	ErrInvalidValue UsingErrorType = iota
	ErrInvalidType
	ErrUnknownTag
	ErrInvalidTag
	ErrFormatTag
	ErrUnsupported
)

type UsingError struct {
	Type  UsingErrorType
	Where string
	Err   error
}

func (v UsingError) Error() string {
	switch v.Type {
	case ErrInvalidValue:
		return fmt.Sprintf("invalid value %s in %s", v.Err.Error(), v.Where)
	case ErrInvalidType:
		return fmt.Sprintf("invalid type %s in %s", v.Err.Error(), v.Where)
	case ErrUnknownTag:
		return fmt.Sprintf("unknown tag %s in %s", v.Err.Error(), v.Where)
	case ErrInvalidTag:
		return fmt.Sprintf("invalid tag %s in %s", v.Err.Error(), v.Where)
	case ErrFormatTag:
		return fmt.Sprintf("format tag %s in %s", v.Err.Error(), v.Where)
	case ErrUnsupported:
		return fmt.Sprintf("unsupported %s in %s", v.Err.Error(), v.Where)
	}
	return fmt.Sprintf("unknown error %s in %s", v.Err.Error(), v.Where)
}
