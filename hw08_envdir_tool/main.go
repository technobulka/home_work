package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println("usage: go-envdir /path/to/env/dir command arg1 arg2")
		os.Exit(ExitCodeError)
	}

	envDir, command := args[1], args[2:]

	env, err := ReadDir(envDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(ExitCodeError)
	}

	result := RunCmd(command, env)
	os.Exit(result)
}
