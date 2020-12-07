package hw02_unpack_string //nolint:golint,stylecheck

import (
	"errors"
	"strconv"
	"strings"
	"unicode"
)

var ErrInvalidString = errors.New("invalid string")

func Unpack(s string) (string, error) {
	var letter string
	var escape bool
	var b strings.Builder

	for _, char := range s {
		if char == 92 && !escape {
			escape = true
			continue
		}

		if !unicode.IsDigit(char) || escape {
			if letter != "" {
				b.WriteString(strings.Repeat(letter, 1))
			}

			if escape && unicode.IsLetter(char) {
				return "", ErrInvalidString
			}

			letter = string(char)
			escape = false
			continue
		}

		if unicode.IsDigit(char) {
			if letter == "" {
				return "", ErrInvalidString
			}

			d, err := strconv.Atoi(string(char))

			if err != nil {
				return "", ErrInvalidString
			}

			b.WriteString(strings.Repeat(letter, d))
			letter = ""
		}
	}

	b.WriteString(strings.Repeat(letter, 1))

	return b.String(), nil
}
