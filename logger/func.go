package logger

import (
	"fmt"
	"os"
	"sync/atomic"
)

func writeln(fdname string, msg string) (int, error) {

	filename := getFilename(fdname)
	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, fileMode)
	defer fd.Close()
	if err != nil {
		return 0, err
	}
	return fmt.Fprintln(fd, msg)
}

func getFile(fdName string) (*os.File, error) {

	fd, ok := fdMaps[fdName]
	if !ok {
		flag := atomic.AddInt32(&mux, 1)
		if flag == 1 {
			go func() {
				err := build(fdName)
				if err != nil {
					writeln(fdName, "error: "+err.Error())
				}
			}()
		}
		filename := getFilename(fdName)
		return os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, fileMode)
	} else {
		return fd, nil
	}
}

func getFilename(file string, suffix ...string) string {

	_suffix := logSuffix
	if len(suffix) > 0 {
		_suffix = suffix[0]
	}
	if len(_suffix) > 0 {
		return logPath + "/" + file + ".log." + _suffix
	}
	return logPath + "/" + file + ".log"

}

func swapSuffix(suffix string) error {

	if len(fdMaps) > 0 {
		newFdMaps := make(map[string]*os.File)
		tmpFdMaps := make(map[string]*os.File)
		for _file, _ := range fdMaps {
			filename := getFilename(_file, suffix)
			fd, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
			if err != nil {
				writeln("logger.error", err.Error())
				continue
			}
			newFdMaps[_file] = fd
		}
		tmpFdMaps = fdMaps
		fdMaps = newFdMaps
		for _, __fd := range tmpFdMaps {
			if err := __fd.Close(); err != nil {
				writeln("logger.error", err.Error())
				continue
			}
		}
		logSuffix = suffix
	}
	atomic.StoreInt32(&mux, 0)
	return nil
}

func build(file string) error {

	filename := getFilename(file)
	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0755)
	if err != nil {
		return err
	}

	newFdMaps := make(map[string]*os.File)

	for _file, _fd := range fdMaps {
		newFdMaps[_file] = _fd
	}
	newFdMaps[file] = fd
	fdMaps = newFdMaps
	atomic.StoreInt32(&mux, 0)
	return nil
}
