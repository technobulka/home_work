package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	in, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer out.Close()

	fileInfo, err := in.Stat()
	if err != nil {
		return err
	}

	if fileInfo.IsDir() {
		return ErrUnsupportedFile
	}

	if fileInfo.Size() < offset {
		return ErrOffsetExceedsFileSize
	}

	_, err = in.Seek(offset, io.SeekStart)
	if err != nil {
		return err
	}

	if limit == 0 {
		limit = fileInfo.Size() - offset
	}

	count := limit
	if fileInfo.Size()-offset < limit {
		count = fileInfo.Size() - offset
	}

	bar := pb.Simple.Start64(count)
	barReader := bar.NewProxyReader(in)

	_, err = io.CopyN(out, barReader, limit)
	if err != nil && err != io.EOF {
		return err
	}

	bar.Finish()

	return nil
}
