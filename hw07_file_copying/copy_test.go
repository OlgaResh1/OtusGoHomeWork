package main

import (
	"crypto/md5"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"testing"
)

type TestCast struct {
	input  string
	result string
	limit  int64
	offset int64
}

func md5Sum(f *os.File) (string, error) {
	file1Sum := md5.New()
	_, err := io.Copy(file1Sum, f)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%X", file1Sum.Sum(nil)), nil
}

func isEqualFiles(fileName1 string, fileName2 string) (bool, error) {
	var sum1, sum2 string
	file1, err := os.Open(fileName1)
	if err != nil {
		return false, err
	}
	defer file1.Close()
	if sum1, err = md5Sum(file1); err != nil {
		return false, err
	}
	file2, err := os.Open(fileName2)
	if err != nil {
		return false, err
	}
	defer file2.Close()
	if sum2, err = md5Sum(file2); err != nil {
		return false, err
	}
	if sum1 == sum2 {
		return true, nil
	}
	return false, nil
}

func TestCopy(t *testing.T) {
	testsGood := [...]TestCast{
		{
			input:  "testdata/input.txt",
			result: "testdata/out_offset0_limit0.txt",
			limit:  0,
			offset: 0,
		},
		{
			input:  "testdata/input.txt",
			result: "testdata/out_offset0_limit10.txt",
			limit:  10,
			offset: 0,
		},
		{
			input:  "testdata/input.txt",
			result: "testdata/out_offset0_limit1000.txt",
			limit:  1000,
			offset: 0,
		},
		{
			input:  "testdata/input.txt",
			result: "testdata/out_offset0_limit10000.txt",
			limit:  10000,
			offset: 0,
		},
		{
			input:  "testdata/input.txt",
			result: "testdata/out_offset100_limit1000.txt",
			limit:  1000,
			offset: 100,
		},
		{
			input:  "testdata/input.txt",
			result: "testdata/out_offset6000_limit1000.txt",
			limit:  1000,
			offset: 6000,
		},
	}

	for _, test := range testsGood {
		tempFileName := filepath.Join(os.TempDir(), "test.txt")
		err := Copy(test.input, tempFileName, test.offset, test.limit)
		if err != nil {
			t.Errorf("Copy failed: %v", err)
			os.Remove(tempFileName)
			continue
		}
		if _, err := os.Stat(tempFileName); os.IsNotExist(err) {
			t.Errorf("Copy failed, file not exist: %v", err)
			continue
		}

		ok, err := isEqualFiles(test.result, tempFileName)
		if err != nil {
			t.Errorf("Error defined md5 file sum %v", err)
		} else if !ok {
			t.Errorf("Copy failed, files not equal: %v", test.result)
		}
		os.Remove(tempFileName)
	}
}

func TestCopyErr(t *testing.T) {
	testsErr := [...]TestCast{
		{
			input:  "testdata/notfound.txt",
			limit:  0,
			offset: 0,
		},
		{
			input:  "testdata/input.txt",
			limit:  -1,
			offset: 100000,
		},
	}
	for _, test := range testsErr {
		tempFileName := filepath.Join(os.TempDir(), "test.txt")
		err := Copy(test.input, tempFileName, test.offset, test.limit)
		if err == nil {
			t.Errorf("Copy should fail: %v", err)
			continue
		}
		os.Remove(tempFileName)
	}
}
