package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

type UserRole string

// Test the function on different structures and other types.
type (
	User struct {
		ID     string `json:"id" validate:"len:32"`
		Name   string
		Age    int      `validate:"min:18|max:50"`
		Email  string   `validate:"regexp:^\\w+@\\w+\\.\\w+$"`
		Role   UserRole `validate:"in:admin,stuff"`
		Phones []string `validate:"len:11"`
		meta   json.RawMessage
	}

	App struct {
		Version string `validate:"regexp:^\\d(\\.\\d){2}$|len:5"`
	}

	Cop struct {
		Good User `validate:"nested"`
		Bad  User `validate:"nested"`
	}

	Token struct {
		Header    []byte
		Payload   []byte
		Signature []byte
	}
)

func TestValidate(t *testing.T) {
	var noErrors ValidationErrors

	tests := []struct {
		name        string
		in          interface{}
		expectedErr error
	}{
		{
			"not struct",
			1,
			fmt.Errorf("not struct"),
		},
		{
			"simple no errors",
			App{"1.0.0"},
			noErrors,
		},
		{
			"simple one error",
			App{"11.0.0"},
			ValidationErrors{
				ValidationError{"App.Version", fmt.Errorf("invalid regexp, invalid length")},
			},
		},
		{
			"all rule types",
			User{
				ID:     "5d41402abc4b2a76b9719d911017c592",
				Name:   "Tester",
				Age:    25,
				Email:  "tester@test.com",
				Role:   "stuff",
				Phones: []string{"01234567890", "99999999999"},
				meta:   nil,
			},
			noErrors,
		},
		{
			"nested validation",
			Cop{
				User{
					ID:     "5d41402abc4b2a76b9719d911017c592",
					Name:   "Good cop",
					Age:    31,
					Email:  "good@test.com",
					Role:   "stuff",
					Phones: []string{"01234567890"},
					meta:   nil,
				},
				User{
					ID:     "---",
					Name:   "Bad cop",
					Age:    55,
					Email:  "bad@com",
					Role:   "cop",
					Phones: []string{"123"},
					meta:   nil,
				},
			},
			ValidationErrors{
				ValidationError{"Cop.Bad.ID", fmt.Errorf("invalid length")},
				ValidationError{"Cop.Bad.Age", fmt.Errorf("invalid max")},
				ValidationError{"Cop.Bad.Email", fmt.Errorf("invalid regexp")},
				ValidationError{"Cop.Bad.Role", fmt.Errorf("invalid contains")},
				ValidationError{"Cop.Bad.Phones", fmt.Errorf("invalid length")},
			},
		},
		{
			"without validation",
			Token{
				[]byte("test"),
				[]byte("no"),
				[]byte("validates"),
			},
			noErrors,
		},
	}

	for _, tt := range tests {
		tt := tt

		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.expectedErr, Validate(tt.in))
		})
	}
}
