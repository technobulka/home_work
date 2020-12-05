package hw02_unpack_string //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

type test struct {
	input    string
	expected string
	err      error
}

func TestUnpack(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    "a4bc2d5e",
			expected: "aaaabccddddde",
		},
		{
			input:    "abccd",
			expected: "abccd",
		},
		{
			input:    "3abc",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "45",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "aaa10b",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "",
			expected: "",
		},
		{
			input:    "aaa0b",
			expected: "aab",
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}

func TestUnpackWithUpper(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    "ASD",
			expected: "ASD",
		},
		{
			input:    "A3D0",
			expected: "AAA",
		},
		{
			input:    "3A",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "SSS3S",
			expected: "SSSSSS",
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}

func TestUnpackWithSpecialChars(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    `a\`,
			expected: "a",
		},
		{
			input:    `\\3`,
			expected: `\\\`,
		},
		{
			input:    "a\n3abc",
			expected: "a\n\n\nabc",
		},
		{
			input:    "a\t2b",
			expected: "a\t\tb",
		},
		{
			input:    "³3",
			expected: "³³³",
		},
		{
			input:    "#16",
			expected: "",
			err:      ErrInvalidString,
		},
		{
			input:    "=8=8",
			expected: "================",
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}

func TestUnpackWithLanguages(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    "山3",
			expected: "山山山",
		},
		{
			input:    "montañ3a0",
			expected: "montañññ",
		},
		{
			input:    "جبل4",
			expected: "جبلللل",
		},
		{
			input:    "ภูเขา2",
			expected: "ภูเขาา",
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}

func TestUnpackWithEscape(t *testing.T) {
	for _, tst := range [...]test{
		{
			input:    `qwe\4\5`,
			expected: `qwe45`,
		},
		{
			input:    `qwe\45`,
			expected: `qwe44444`,
		},
		{
			input:    `qwe\\5`,
			expected: `qwe\\\\\`,
		},
		{
			input:    `qwe\\\3`,
			expected: `qwe\3`,
		},
		{
			input:    `qw\ne`,
			expected: "",
			err:      ErrInvalidString,
		},
	} {
		result, err := Unpack(tst.input)
		require.Equal(t, tst.err, err)
		require.Equal(t, tst.expected, result)
	}
}
