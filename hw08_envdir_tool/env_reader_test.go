package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestErrorFindDir(t *testing.T) {
	_, err := ReadDir("./testdata/notexist")

	require.Error(t, err)
}

func TestReadDir(t *testing.T) {
	res, err := ReadDir("./testdata/env")
	require.NoError(t, err)

	if len(res) != 5 {
		require.Equal(t, 5, len(res), "Expected 2 variables, got %d")
	}

	require.Equal(t, res["BAR"].Value, "bar")
	require.Equal(t, res["BAR"].NeedRemove, false)

	require.Equal(t, res["EMPTY"].Value, "")
	require.Equal(t, res["EMPTY"].NeedRemove, false)

	require.Equal(t, res["FOO"].Value, "   foo\nwith new line")
	require.Equal(t, res["FOO"].NeedRemove, false)

	require.Equal(t, res["HELLO"].Value, `"hello"`)
	require.Equal(t, res["HELLO"].NeedRemove, false)

	require.Equal(t, res["UNSET"].Value, "")
	require.Equal(t, res["UNSET"].NeedRemove, true)
}
