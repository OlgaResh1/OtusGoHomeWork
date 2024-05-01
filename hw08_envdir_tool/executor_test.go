package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

type Command struct {
	testName        string
	cmd             []string
	expectedEnv     Environment
	expectedRetCode int
}

func TestRunCmd(t *testing.T) {
	testCases := []Command{
		{
			testName: "Test env variables",
			cmd:      []string{"ls", "-l"},
			expectedEnv: Environment{
				"HOME": {"/home/user", false},
				"PATH": {"/usr/bin", false},
			},
			expectedRetCode: 0,
		},
		{
			testName:        "Echo first argument",
			cmd:             []string{"echo", "$0"},
			expectedRetCode: 0,
		},
		{
			testName:        "Echo hello",
			cmd:             []string{"echo", "$UNAME"},
			expectedEnv:     Environment{"UNAME": {"echo", false}},
			expectedRetCode: 0,
		},
		{
			testName:        "Command not found",
			cmd:             []string{"/bin/bash", "bad_command", "-v"},
			expectedRetCode: 127,
		},
		{
			testName:        "No such file or directory",
			cmd:             []string{"/bin/bash", "./bad_command", "-v"},
			expectedRetCode: 127,
		},
	}
	for _, testCase := range testCases {

		returnCode := RunCmd(testCase.cmd, testCase.expectedEnv)
		require.Equalf(t, testCase.expectedRetCode, returnCode, testCase.testName)

		for envKey, envVal := range testCase.expectedEnv {
			if !envVal.NeedRemove {
				require.Equalf(t, os.Getenv(envKey), envVal.Value, "Expected environment variable %s was not valid", envVal.Value)
			}
		}
	}
}
