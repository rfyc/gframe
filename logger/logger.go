package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"

	// "runtime"
	"sync/atomic"
	"time"

	"github.com/phper-go/frame/func/conv"
)

func SetSuffix(suffix string) {

	if suffix == "" {
		suffix = time.Now().Format("20060102")
	}
	if logSuffix != suffix {
		flag := atomic.AddInt32(&mux, 1)
		if flag == 1 {
			go func() {
				err := swapSuffix(suffix)
				if err != nil {
					writeln("logger.error", err.Error())
				}
			}()
		}
	}
}

func SetPath(path string) (realpath string, err error) {

	//******* args path ********/
	if path != "" {
		_, err = os.Stat(path)
		if err != nil {
			err = os.MkdirAll(path, fileMode)
		}
		if err == nil {
			logPath = path
			return logPath, nil
		}
	}
	//******* default path ********/
	var ferr error
	if realpath, ferr = filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		realpath += "/logs"
		_, ferr = os.Stat(realpath)
		if ferr != nil {
			ferr = os.MkdirAll(path, fileMode)
		}
		if ferr == nil {
			logPath = realpath
			return logPath, err
		}
	}
	//******* return error ********/
	var returnErr error
	returnErr = errors.New(realpath + ":" + ferr.Error())
	if err != nil {
		returnErr = errors.New(path + ":" + err.Error() + ";" + returnErr.Error())
	}
	return logPath, returnErr
}

func Writeln(fdname string, msg string, noSuffix ...bool) (int, error) {

	nosuf := false
	if len(noSuffix) > 0 {
		nosuf = noSuffix[0]
	}
	if !nosuf {

		fd, ok := fdMaps[fdname]
		if !ok {
			flag := atomic.AddInt32(&mux, 1)
			if flag == 1 {
				go func() {
					err := build(fdname)
					if err != nil {
						writeln("logger.error", err.Error())
					}
				}()
			}
			return writeln(fdname, msg)
		} else {
			len, err := fmt.Fprintln(fd, msg)
			if err != nil {
				return writeln(fdname, msg)
			}
			return len, err
		}
	} else {
		return writeln(fdname, msg)
	}
}

func Format(data ...interface{}) *logFormat {
	logFmt := &logFormat{}
	logFmt.data = data
	logFmt.caller = true
	return logFmt
}

func ErrorHandler(fdname, msg string, size int, err error) {

}

func FormatHandler(caller bool, data ...interface{}) string {

	logmsg := time.Now().Format("2006-01-02 15:04:05.000") + " | "
	logmsg += conv.String(os.Getpid()) + " | "
	if caller {
		// _, file, line, ok := runtime.Caller(2)
		// if ok {
		// 	logmsg += filepath.Base(file) + ":" + conv.String(line) + " | "
		// }
	}
	if len(data) > 0 {
		for _, val := range data {
			logmsg += conv.String(val) + " "
		}
	}
	return logmsg
}

func Access() accessInterface {

	log := &access{}
	log.beginTime = time.Now()
	return log
}
