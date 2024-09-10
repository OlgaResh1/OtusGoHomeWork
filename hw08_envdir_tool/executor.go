package main

import (
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	cmdExec := exec.Command(cmd[0], cmd[1:]...)

	cmdExec.Stdin = os.Stdin
	cmdExec.Stdout = os.Stdout
	cmdExec.Stderr = os.Stderr

	for envName, envVal := range env {
		if envVal.NeedRemove {
			os.Unsetenv(envName)
		} else {
			os.Setenv(envName, envVal.Value)
		}
	}
	cmdExec.Env = os.Environ()
	cmdExec.Run()
	// if err != nil {
	returnCode = cmdExec.ProcessState.ExitCode()
	//}
	return
}
