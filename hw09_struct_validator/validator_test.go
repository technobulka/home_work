package hw09structvalidator

import (
	"encoding/json"
	"fmt"
	"testing"
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
		Version string `validate:"len:5"`
	}

	Post struct {
		User User
		App  App
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
)

func TestValidate(t *testing.T) {
	tests := []struct {
		in          interface{}
		expectedErr error
	}{
		{
			App{"1.0.0"},
			nil,
		},
		{
			App{"11.0.0"},
			fmt.Errorf("not valid"),
		},
		{
			User{
				ID:     "5d41402abc4b2a76b9719d911017c592",
				Name:   "Tester",
				Age:    25,
				Email:  "tester@test.com",
				Role:   "stuff",
				Phones: []string{"01234567890"},
				meta:   nil,
			},
			nil,
		},
		{
			Post{
				User{
					ID:     "5d41402abc4b2a76b9719d911017c592",
					Name:   "Tester",
					Age:    25,
					Email:  "tester@test.com",
					Role:   "stuff",
					Phones: []string{"01234567890"},
					meta:   nil,
				},
				App{"1.0.0"},
			},
			nil,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case %d", i), func(t *testing.T) {
			tt := tt
			t.Parallel()

			_ = Validate(tt.in)
			_ = tt
		})
	}
}
