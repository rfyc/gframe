package element

import (
	"bytes"
	"html/template"
)

func Widget(widgetPath, widgetFile, input_desc string) Interface {
	widget := &widgetElement{}
	widget.itype = "widget"
	widget.desc = input_desc
	widget.options = make(map[string]string)
	widget.path = widgetPath
	widget.file = widgetFile
	return widget
}

type widgetElement struct {
	inputElement
	widgetPath string
	path       string
	file       string
	form       interface{}
}

func (this *widgetElement) Set(formObj interface{}, required bool) {
	this.form = formObj
	this.required = required
}
func (this *widgetElement) Render(attributes ...string) template.HTML {

	render, err := template.ParseFiles(this.path + "/" + this.file + ".html")

	if err != nil {
		return template.HTML("widget " + this.file + " " + err.Error())
	}
	buffer := bytes.NewBuffer([]byte{})
	err = render.Execute(buffer, this.form)
	if err != nil {
		return template.HTML("widget " + this.file + " " + err.Error())
	}
	return template.HTML(buffer.Bytes())
}
