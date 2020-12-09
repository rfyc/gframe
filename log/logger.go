package log

import (
	"fmt"
	"io"
	"os"
)

var Path string

func Format(data ...interface{}) *logFormat {
	logFmt := &logFormat{}
	logFmt.logmsg = LogFormatHandler(data...)
	return logFmt
}

func LogFileHandler() {

}

func ErrorHandler() {

}

func LogFormatHandler(data ...interface{}) string {

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

func writeln(fdname string, msg string) (int, error) {

	filename := getFilename(fdname)
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

	size, err := writeln(fdname, logmsg)
	ErrorHandler(fdname, logmsg, size, err)
	return this
}

func (this *logFormat) Error(fdname string) *logFormat {

	fdname += ".error"
	return this.Writeln(fdname)
}

func (this *logFormat) Echo(stdio io.Writer) *logFormat {
	size, err := fmt.Fprintln(stdio, this.logmsg)
	ErrorHandler("Stdout", logmsg, size, err)
	return this
}
