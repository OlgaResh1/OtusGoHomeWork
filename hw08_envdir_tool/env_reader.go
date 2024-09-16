package main

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

var ErrWrongFileName = errors.New("wrong filename")

func trimText(text []byte) []byte {
	i := bytes.Index(text, []byte{'\n'})
	if i >= 0 {
		text = text[:i]
	}

	text = bytes.TrimRight(text, "\t")
	text = bytes.TrimRight(text, " ")
	text = bytes.ReplaceAll(text, []byte{0}, []byte{'\n'})
	return text
}

func readEnvFromFile(fileName string) (EnvValue, error) {
	item := EnvValue{}

	file, err := os.Open(fileName)
	if err != nil {
		return item, err
	}
	defer file.Close()

	reader := bufio.NewReader(file)
	line, err := reader.ReadString('\n')
	if err != nil && !errors.Is(err, io.EOF) {
		return item, err
	}

	if len(line) == 0 {
		item.NeedRemove = true
		return item, nil
	}

	item.Value = string(trimText([]byte(line)))
	return item, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envs := make(Environment)

	var filePaths []string
	err := filepath.WalkDir(dir,
		func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			if !d.IsDir() {
				filePaths = append(filePaths, path)
			}
			return nil
		})
	if err != nil {
		return nil, err
	}

	for _, fileName := range filePaths {
		env, err := readEnvFromFile(fileName)
		if err != nil {
			return nil, err
		}
		envs[filepath.Base(fileName)] = env
	}
	return envs, nil
}
