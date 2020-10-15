package logger

import (
	"fmt"
	"io"
	"os"
)

type logFormat struct {
	caller bool
	data   []interface{}
}

func (this *logFormat) String() string {
	return FormatHandler(this.caller, this.data...)
}

func (this *logFormat) Info(fdname string) *logFormat {

	fdname += ".info"
	logmsg := FormatHandler(this.caller, this.data...)
	size, err := Writeln(fdname, logmsg)
	ErrorHandler(fdname, logmsg, size, err)
	return this
}

func (this *logFormat) Writeln(fdname string) *logFormat {

	logmsg := FormatHandler(this.caller, this.data...)
	size, err := Writeln(fdname, logmsg)
	ErrorHandler(fdname, logmsg, size, err)
	return this
}

func (this *logFormat) Error(fdname string) *logFormat {

	fdname += ".error"
	logmsg := FormatHandler(this.caller, this.data...)
	size, err := Writeln(fdname, logmsg)
	ErrorHandler(fdname, logmsg, size, err)
	return this
}

func (this *logFormat) Println() *logFormat {
	logmsg := FormatHandler(this.caller, this.data...)
	size, err := fmt.Fprintln(os.Stderr, logmsg)
	ErrorHandler("Stdout", logmsg, size, err)
	return this
}

func (this *logFormat) Echo(stdio io.Writer) *logFormat {
	logmsg := FormatHandler(this.caller, this.data...)
	size, err := fmt.Fprintln(stdio, logmsg)
	ErrorHandler("Stdout", logmsg, size, err)
	return this
}
