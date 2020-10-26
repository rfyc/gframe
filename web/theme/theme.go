package theme

import (
	"bytes"
	"html/template"
	"os"
	"path/filepath"
)

var (
	DefaultPath    string
	DefaultStyle   string
	DefaultSuffix  string
	DefaultLayouts []string
)

type Interface interface {
	Construct()
	SetPath(path string)
	SetStyle(style string)
	SetSuffix(suffix string)
	SetLayouts(layoutFiles []string)
	Assign(key string, val interface{})
	Render(tplFile string, layoutFiles ...string) []byte
}

func init() {

	if DefaultPath == "" {
		rootDir, _ := filepath.Abs(filepath.Dir(os.Args[0]))
		DefaultPath = rootDir + "/themes"
	}
	if DefaultStyle == "" {
		DefaultStyle = "default"
	}

	if DefaultSuffix == "" {
		DefaultSuffix = ".html"
	}
}

type Theme struct {
	path    string
	style   string
	suffix  string
	layouts []string
	params  map[string]interface{}
}

func (this *Theme) Construct() {

	this.params = make(map[string]interface{})
	this.path = DefaultPath
	this.style = DefaultStyle
	this.suffix = DefaultSuffix
}
func (this *Theme) SetPath(path string) {
	this.path = path
}
func (this *Theme) SetStyle(style string) {
	this.style = style
}
func (this *Theme) SetSuffix(suffix string) {
	this.suffix = suffix
}
func (this *Theme) SetLayouts(layoutFiles []string) {
	this.layouts = layoutFiles
}
func (this *Theme) Assign(key string, val interface{}) {
	this.params[key] = val
}
func (this *Theme) Render(tplFile string, layoutFiles ...string) []byte {

	if len(layoutFiles) == 0 {
		layoutFiles = this.layouts
	}
	layoutFiles = append(layoutFiles, tplFile)
	for index, file := range layoutFiles {
		layoutFiles[index] = this.path + "/" + this.style + "/" + file + this.suffix
	}

	render, err := template.ParseFiles(layoutFiles...)
	if err != nil {
		panic(err.Error())
	}

	buffer := bytes.NewBuffer([]byte{})
	err = render.Execute(buffer, this.params)
	if err != nil {
		panic(err.Error())
	}

	return buffer.Bytes()
}
