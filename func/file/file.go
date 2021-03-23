package file

import (
	"bufio"
	"os"
	"path/filepath"
)

func ReadLine(buf *bufio.Reader) (string, error) {
	var line []byte
	for {

		ln, more, err := buf.ReadLine()
		if err != nil {
			return "", err
		}
		if line == nil && !more {
			return string(ln), nil
		}
		line = append(line, ln...)
		if !more {
			break
		}
	}
	return string(line), nil
}

func IsDir(path string) bool {

	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func IsFile(path string) bool {

	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}

		return false
	}
	return true
}

func BinDir() string {
	exec, _ := os.Executable()
	return filepath.Dir(exec)
}
