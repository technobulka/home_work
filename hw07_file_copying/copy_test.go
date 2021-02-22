package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	t.Run("file not found", func(t *testing.T) {
		from := "not_existing_file.txt"
		to := "out.txt"

		err := Copy(from, to, 0, 0)
		os.Remove(to)
		require.Error(t, err)
	})

	t.Run("incorrect input", func(t *testing.T) {
		from := "testdata"
		to := "out.txt"

		err := Copy(from, to, 0, 0)
		os.Remove(to)
		require.Error(t, err)
	})

	t.Run("out of file size", func(t *testing.T) {
		from := "testdata/input.txt"
		to := "out.txt"

		err := Copy(from, to, 9000, 0)
		os.Remove(to)
		require.Error(t, err)
	})
}
