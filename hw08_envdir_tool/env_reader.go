package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	envs := make(Environment, len(files))
	for _, file := range files {
		if file.IsDir() {
			return nil, fmt.Errorf("file %q is a directory", file.Name())
		}

		if strings.Contains(file.Name(), "=") {
			return nil, fmt.Errorf("file %q contain %q symbol", file.Name(), "=")
		}

		env, err := getEnv(filepath.Join(dir, file.Name()))
		if err != nil {
			return nil, err
		}

		envs[file.Name()] = env
	}

	return envs, nil
}

func getEnv(path string) (EnvValue, error) {
	env := EnvValue{}

	file, err := os.Open(path)
	if err != nil {
		return env, err
	}
	defer file.Close()

	info, err := file.Stat()
	if err != nil {
		return env, err
	}

	if info.Size() == 0 {
		env.NeedRemove = true
		return env, nil
	}

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	val := strings.TrimRight(scanner.Text(), " \t")
	val = strings.ReplaceAll(val, string([]byte{'\x00'}), "\n")
	env.Value = val

	return env, nil
}
