package logger

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"github.com/phper-go/frame/func/conv"
)

var Path string
var LogFormatHandler func(data ...interface{}) string
var LogFileHandler func(fdname string) string
var ErrorHandler func(fdname, logmsg string, size int, err error)

func init() {
	Path = filepath.Dir(os.Args[0]) + "/logs"
	LogFormatHandler = logFormatHandler
	LogFileHandler = logFileHandler
	ErrorHandler = errorHandler
}

func Format(data ...interface{}) *logFormat {
	logFmt := &logFormat{}
	logFmt.logmsg = LogFormatHandler(data...)
	return logFmt
}

func logFileHandler(fdname string) string {

	return Path + "/" + fdname + ".log." + time.Now().Format("20060102")
}

func errorHandler(fdname, logmsg string, size int, err error) {

}

func logFormatHandler(data ...interface{}) string {

	logmsg := time.Now().Format("2006-01-02 15:04:05.000") + " | "
	if len(data) > 0 {
		for _, val := range data {
			logmsg += conv.String(val) + " "
		}
	}
	return logmsg
}

func writeln(fdname string, msg string) (int, error) {

	var filename = LogFileHandler(fdname)
	var fileMode os.FileMode = 0755
	fd, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, fileMode)
	defer fd.Close()
	if err != nil {
		return 0, err
	}
	return fmt.Fprintln(fd, msg)
}

type logFormat struct {
	logmsg string
}

func (this *logFormat) Info(fdname string) *logFormat {

	fdname += ".info"
	return this.Writeln(fdname)
}

func (this *logFormat) Writeln(fdname string) *logFormat {

	size, err := writeln(fdname, this.logmsg)
	ErrorHandler(fdname, this.logmsg, size, err)
	return this
}

func (this *logFormat) Error(fdname string) *logFormat {

	fdname += ".error"
	return this.Writeln(fdname)
}

func (this *logFormat) Echo(stdio io.Writer) *logFormat {
	size, err := fmt.Fprintln(stdio, this.logmsg)
	ErrorHandler("Writer", this.logmsg, size, err)
	return this
}
