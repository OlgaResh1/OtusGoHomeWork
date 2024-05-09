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

	require.Len(t, res, 5, "Expected 5 variables, got %d")

	require.Equal(t, "bar", res["BAR"].Value)
	require.Equal(t, false, res["BAR"].NeedRemove)

	require.Equal(t, "", res["EMPTY"].Value)
	require.Equal(t, false, res["EMPTY"].NeedRemove)

	require.Equal(t, "   foo\nwith new line", res["FOO"].Value)
	require.Equal(t, false, res["FOO"].NeedRemove)

	require.Equal(t, `"hello"`, res["HELLO"].Value)
	require.Equal(t, false, res["HELLO"].NeedRemove)

	require.Equal(t, "", res["UNSET"].Value)
	require.Equal(t, true, res["UNSET"].NeedRemove)
}
