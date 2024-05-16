package hw09structvalidator

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}

	App struct {
		Version string `validate:"len:5"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}

	Response struct {
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty"`
	}
	CardType struct {
		Suit  []string `validate:"in:diamonds,hearts,clubs,spades"`
		Value []int    `validate:"min:2|max:15"`
	}
	CardBadType struct { // wrong len tag for int field
		Suit  []string `validate:"in:diamonds,hearts,clubs,spades"`
		Value []int    `validate:"len:2"`
	}
	UserBad struct { // wrong regexp
		ID     string `json:"id" validate:"len:36"`
		Name   string
		Age    int             `validate:"min:18|max:50"`
		Email  string          `validate:"regexp:^\\qw+@\\w+\\.\\w+$"`
		Role   UserRole        `validate:"in:admin,stuff"`
		Phones []string        `validate:"len:11"`
		meta   json.RawMessage //nolint:unused
	}
	ResponseBad struct { // min tag for string field
		Code int    `validate:"in:200,404,500"`
		Body string `json:"omitempty" validate:"min:100"`
	}
	AppBad struct { // wrong tag name
		Version string `validate:"stringlen:5"`
	}
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: CardType{
				Suit:  []string{"hearts", "hearts"},
				Value: []int{6, 5},
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:    "db85484b-0720-4c5da00e1-a433d5b2944a",
				Name:  "Ivan",
				Age:   35,
				Email: "ivan@example.com",
				Role:  "stuff",
			},
			expectedErr: nil,
		},
		{
			in: User{
				ID:    "db85484b-0720-4c5da00e-a433d5b2944",
				Name:  "A",
				Age:   1,
				Email: "@example.com",
				Role:  "stuff",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "ID", Err: fmt.Errorf("len does not match 34")},
				ValidationError{Field: "Age", Err: fmt.Errorf("value is out of range")},
				ValidationError{Field: "Email", Err: fmt.Errorf("value must match regexp")},
			},
		},
		{
			in: User{
				ID:    "db85484b-0720-4c5da00e1-a433d5b2944a",
				Name:  "B",
				Age:   30,
				Email: "b@example.com",
				Role:  "manager",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Role", Err: fmt.Errorf("value must be in list")},
			},
		},
		{
			in: App{
				Version: "1.0",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Version", Err: fmt.Errorf("len does not match 3")},
			},
		},
		{
			in: Response{
				Code: 400,
				Body: "",
			},
			expectedErr: ValidationErrors{
				ValidationError{Field: "Code", Err: fmt.Errorf("value must be in list")},
			},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			err := Validate(tt.in)
			compareValidateErrors(t, err, tt.expectedErr)
		})
	}
}

func compareValidateErrors(t *testing.T, err error, expected error) {
	t.Helper()
	switch {
	case err == nil:
		if expected != nil {
			t.Errorf("expected error %v, got nil", expected)
		}
	case !errors.As(err, &ValidationErrors{}):
		t.Errorf("expected error %v, got nil", expected)
	default:
		if expected == nil {
			t.Errorf("expected nil error, got %v", err)
		}
		if expected.Error() != err.Error() {
			t.Errorf("expected error %v not valid", err)
		}
	}
}

func TestUsingErrors(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			in: CardBadType{
				Suit:  []string{"hearts", "hearts"},
				Value: []int{6, 5},
			},
			expectedErr: UsingError{Type: ErrUnsupported, Where: "Value", Err: fmt.Errorf("unsupported validate type: []int")},
		},
		{
			in: UserBad{
				ID:    "db85484b-0720-4c5da00e1-a433d5b2944a",
				Name:  "Ivan",
				Age:   35,
				Email: "ivan@example.com",
				Role:  "stuff",
			},
			expectedErr: UsingError{Type: ErrFormatTag, Where: "Email", Err: fmt.Errorf("error compile regexp")},
		},
		{
			in: ResponseBad{
				Code: 404,
				Body: "<html/>",
			},
			expectedErr: UsingError{Type: ErrUnsupported, Where: "Body", Err: fmt.Errorf("unsupported validate type: string")},
		},
		{
			in: AppBad{
				Version: "1.0.0",
			},
			expectedErr: UsingError{Type: ErrInvalidTag, Where: "Version", Err: fmt.Errorf("validate tag not exist")},
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()
			err := Validate(tt.in)
			switch {
			case err == nil:
				if tt.expectedErr != nil {
					t.Errorf("expected error %v, got nil", tt.expectedErr)
				}
			case !errors.As(err, &UsingError{}):
				t.Errorf("expected error %v, got nil", tt.expectedErr)
			default:
				if tt.expectedErr == nil {
					t.Errorf("expected nil error, got %v", err)
				}
				if err.Error() != tt.expectedErr.Error() {
					t.Errorf("expected error %v ", tt.expectedErr)
				}
			}
		})
	}
}
