package main

import (
	"errors"
	"log"
	"os"
	"os/exec"
)

const (
	ExitCodeOk    = 0
	ExitCodeError = 1
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	for key, val := range env {
		if err := os.Unsetenv(key); err != nil {
			log.Fatal(err)
			return ExitCodeError
		}

		if !val.NeedRemove {
			if err := os.Setenv(key, val.Value); err != nil {
				log.Fatal(err)
				return ExitCodeError
			}
		}
	}

	if len(cmd) == 0 {
		return ExitCodeError
	}

	var args []string
	if len(cmd) > 1 {
		args = cmd[1:]
	}

	c := exec.Command(cmd[0], args...) //nolint:gosec
	c.Stdin = os.Stdin
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr

	if err := c.Run(); err != nil {
		var exitError *exec.ExitError
		if errors.As(err, &exitError) {
			return exitError.ExitCode()
		}

		return ExitCodeError
	}

	return ExitCodeOk
}
