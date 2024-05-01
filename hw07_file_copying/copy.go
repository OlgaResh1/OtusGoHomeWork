package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

type fileCopier struct {
	input  *os.File
	output *os.File
	limit  int64
	offset int64
}

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrDefFileSize           = errors.New("error defining file size")
)

func (c *fileCopier) copyWithStatusBar() error {
	bar := pb.Full.Start64(c.limit)
	defer bar.Finish()
	barReader := bar.NewProxyReader(c.input)
	var err error
	if c.limit > 0 {
		_, err = io.CopyN(c.output, barReader, c.limit)
	} else {
		_, err = io.Copy(c.output, barReader)
	}
	if errors.Is(err, io.EOF) {
		return nil
	}
	return err
}

func (c *fileCopier) initReader(fromPath string) error {
	fileFrom, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	fi, err := fileFrom.Stat()
	if err != nil {
		return ErrUnsupportedFile
	}

	if c.offset > fi.Size() {
		return ErrOffsetExceedsFileSize
	}

	fileFrom.Seek(c.offset, 0)

	c.input = fileFrom
	return nil
}

func (c *fileCopier) initWriter(toPath string) error {
	toFile, err := os.OpenFile(toPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0o644)
	if err != nil {
		return err
	}
	c.output = toFile
	return nil
}

func (c *fileCopier) close() error {
	err := c.input.Close()
	if err != nil {
		return err
	}
	err = c.output.Close()
	if err != nil {
		return err
	}
	return nil
}

func Copy(fromPath, toPath string, offset, limit int64) error {
	c := fileCopier{offset: offset, limit: limit}
	err := c.initReader(fromPath)
	if err != nil {
		return err
	}
	defer c.close()

	err = c.initWriter(toPath)
	if err != nil {
		return err
	}

	err = c.copyWithStatusBar()
	if err != nil {
		return err
	}

	return err
}
