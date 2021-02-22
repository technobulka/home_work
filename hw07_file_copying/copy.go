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

func getSize(file *os.File, offset, limit int64) (int64, error) {
	fileInfo, err := file.Stat()
	if err != nil {
		return 0, err
	}

	if fileInfo.IsDir() {
		return 0, ErrUnsupportedFile
	}

	if fileInfo.Size() < offset {
		return 0, ErrOffsetExceedsFileSize
	}

	_, err = file.Seek(offset, io.SeekStart)
	if err != nil {
		return 0, err
	}

	return fileInfo.Size(), nil
}

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

	size, err := getSize(in, offset, limit)
	if err != nil {
		return err
	}

	if limit == 0 {
		limit = size - offset
	}

	count := limit
	if size-offset < limit {
		count = size - offset
	}

	bar := pb.Simple.Start64(count)
	barReader := bar.NewProxyReader(in)

	_, err = io.CopyN(out, barReader, limit)
	if err != nil && !errors.Is(err, io.EOF) {
		return err
	}

	bar.Finish()

	return nil
}
