package main

import (
	"bytes"
	"errors"
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

var ErrWrongFileName = errors.New("wrong filename")

func fileNamesInDir(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var fnames []string
	for _, file := range files {
		if file.IsDir() {
			subFiles, err := fileNamesInDir(filepath.Join(dir, file.Name()))
			if err != nil {
				return nil, err
			}
			fnames = append(fnames, subFiles...)
		} else {
			if strings.Contains(file.Name(), "=") {
				return nil, ErrWrongFileName
			}
			fnames = append(fnames, filepath.Join(dir, file.Name()))
		}
	}

	return fnames, nil
}

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

	fbytes, err := os.ReadFile(fileName)
	if err != nil {
		return item, err
	}

	if len(fbytes) == 0 {
		item.NeedRemove = true
		return item, nil
	}

	item.Value = string(trimText(fbytes))
	return item, nil
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	envs := make(Environment)

	fileNames, err := fileNamesInDir(dir)
	if err != nil {
		return nil, err
	}

	for _, fileName := range fileNames {
		env, err := readEnvFromFile(fileName)
		if err != nil {
			return nil, err
		}
		envs[filepath.Base(fileName)] = env
	}
	return envs, nil
}
