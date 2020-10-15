package logger

import (
	"os"
	"path/filepath"
	"time"
)

var logSuffix string
var logPath string
var fileMode os.FileMode = 0755
var fdMaps map[string]*os.File
var mux int32
var fdRunlog *os.File
var guid string

func init() {

	SetPath(filepath.Dir(os.Args[0]) + "/logs")
	logSuffix = time.Now().Format("20060102")
	fdMaps = make(map[string]*os.File)
}
