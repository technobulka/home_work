package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testPath = "test_dir"
)

func TestReadDir(t *testing.T) {
	t.Run("Empty dir", func(t *testing.T) {
		t.Cleanup(removeTestDir)

		if err := os.Mkdir(testPath, 0777); err != nil {
			fmt.Println(err)
		}

		env, err := ReadDir(testPath)
		require.NoError(t, err)
		require.Len(t, env, 0)
	})

	t.Run("Dir in envs", func(t *testing.T) {
		t.Cleanup(removeTestDir)

		if err := os.Mkdir(testPath, 0755); err != nil {
			fmt.Println(err)
		}

		if err := os.Mkdir(path.Join(testPath, "ENV"), 0755); err != nil {
			fmt.Println(err)
		}

		env, err := ReadDir(testPath)
		require.Error(t, err)
		require.Len(t, env, 0)
	})

	t.Run("Special char in file name", func(t *testing.T) {
		t.Cleanup(removeTestDir)

		if err := os.Mkdir(testPath, 0755); err != nil {
			fmt.Println(err)
		}

		file, err := os.Create(path.Join(testPath, "SOME=ENV"))
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()

		env, err := ReadDir(testPath)
		require.Error(t, err)
		require.Len(t, env, 0)
	})

	t.Run("Success envs", func(t *testing.T) {
		t.Cleanup(removeTestDir)

		if err := os.Mkdir(testPath, 0755); err != nil {
			fmt.Println(err)
		}

		files := []struct {
			name string
			data []byte
		}{
			{"FOO", []byte("foo\nsecond line")},
			{"EMPTY", []byte(" \n")},
			{"BAR", []byte("bar\x00new line")},
			{"HELLO", []byte(`"hello"`)},
			{"UNSET", nil},
		}

		for _, file := range files {
			err := ioutil.WriteFile(path.Join(testPath, file.name), file.data, 0644)
			if err != nil {
				fmt.Println(err)
			}
		}

		expect := Environment{
			"FOO": EnvValue{
				Value: "foo",
			},
			"EMPTY": EnvValue{
				Value: "",
			},
			"BAR": EnvValue{
				Value: "bar\nnew line",
			},
			"HELLO": EnvValue{
				Value: `"hello"`,
			},
			"UNSET": EnvValue{
				NeedRemove: true,
			},
		}

		env, err := ReadDir(testPath)
		require.NoError(t, err)
		require.Equal(t, expect, env)
	})
}

func removeTestDir() {
	if err := os.RemoveAll(testPath); err != nil {
		fmt.Println(err)
	}
}
