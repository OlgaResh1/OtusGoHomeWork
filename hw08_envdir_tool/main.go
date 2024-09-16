package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args
	if len(args) < 3 {
		fmt.Println("Usage:	go-envdir pathToDir process arguments ")
		os.Exit(1)
	}
	env, err := ReadDir(args[1])
	if err != nil {
		fmt.Println(err)
		os.Exit(2)
	}
	returnCode := RunCmd(args[2:], env)
	os.Exit(returnCode)
}
