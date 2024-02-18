package hw02unpackstring

import (
	"errors"
	"strings"
)

var ErrInvalidString = errors.New("invalid string")

const invalidSymbol rune = -1

func isDigit(sym rune) bool {
	return (sym >= '0' && sym <= '9') // strconv.Atoi(string(sym))?
}

func isSlash(sym rune) bool {
	return sym == '\\'
}

func Unpack(input string) (string, error) {
	result := strings.Builder{}
	lastSym := invalidSymbol
	slash := false

	for _, sym := range input {
		if slash {
			if !isDigit(sym) && !isSlash(sym) {
				return "", ErrInvalidString
			}
			slash = false
			lastSym = sym
			continue
		}

		if isDigit(sym) {
			if lastSym == invalidSymbol {
				return "", ErrInvalidString
			}
			result.WriteString(strings.Repeat(string(lastSym), int(sym-'0')))
			lastSym = invalidSymbol
			continue
		}

		if isSlash(sym) {
			slash = true
		}

		if lastSym != invalidSymbol {
			result.WriteRune(lastSym)
		}
		lastSym = sym
	}
	if lastSym != invalidSymbol {
		result.WriteRune(lastSym)
	}
	return result.String(), nil
}
