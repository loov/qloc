package main

import (
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CountLines(path string) (*Count, error) {
	count := &Count{
		Ext:   strings.ToLower(filepath.Ext(path)),
		Files: 1,
	}

	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

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
			case 0x0, 0x1, 0x2, 0x3, 0x4, 0x5, 0x6, 0x7,
				0x8, 0x0B, 0x0C, 0x0E, 0x0F, 0x10, 0x11, 0x12, 0x13, 0x14,
				0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f:
				return nil, ErrBinaryFile
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

	return count, nil
}
