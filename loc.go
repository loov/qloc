package main

import (
	"errors"
	"io"
	"os"
	"path/filepath"
	"strings"
)

var (
	ErrEmptyFile = errors.New("empty file")
)

func CountLines(path string) (*Count, error) {
	count := &Count{
		Ext:   strings.ToLower(strings.TrimPrefix(filepath.Ext(path), ".")),
		Files: 1,
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() { _ = file.Close() }()

	stat, err := file.Stat()
	if err != nil {
		return nil, err
	}
	if stat.Size() <= 0 {
		return nil, ErrEmptyFile
	}

	buf := make([]byte, 8196)
	emptyline := true
	for {
		n, err := file.Read(buf)
		if err != nil && err != io.EOF {
			return count, err
		}

		for _, c := range buf[:n] {
			switch c {
			case 0x0:

				count.Blank = 0
				count.Code = 0
				count.Files = 0
				count.Binary = 1

				return count, nil
			case '\n':
				if emptyline {
					count.Blank++
				} else {
					count.Code++
				}
				emptyline = true
			case '\r', ' ', '\t': // ignore
			default:
				emptyline = false
			}
		}
		if err == io.EOF {
			break
		}
	}

	if !emptyline {
		count.Code++
	}

	return count, nil
}
