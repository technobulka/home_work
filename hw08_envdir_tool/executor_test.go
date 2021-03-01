package main

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"testing"
)

const (
	testFile = "test_file.sh"
)

func TestRunCmd(t *testing.T) {
	t.Run("no command", func(t *testing.T) {
		code := RunCmd([]string{}, Environment{})
		require.Equal(t, ExitCodeError, code)
	})

	t.Run("command not found", func(t *testing.T) {
		createTestFile("#!/usr/bin/env bash\nqwe")
		t.Cleanup(removeTestFile)
		args := []string{"bash", testFile}

		code := RunCmd(args, Environment{})
		require.Equal(t, 127, code)
	})

	t.Run("no args", func(t *testing.T) {
		createTestFile("#!/usr/bin/env bash\ngrep")
		t.Cleanup(removeTestFile)
		args := []string{"bash", testFile}

		code := RunCmd(args, Environment{})
		require.Equal(t, 2, code)
	})
}

func createTestFile(code string) {
	err := ioutil.WriteFile(testFile, []byte(code), 0644)
	if err != nil {
		fmt.Println(err)
	}
}

func removeTestFile() {
	if err := os.Remove(testFile); err != nil {
		fmt.Println(err)
	}
}
