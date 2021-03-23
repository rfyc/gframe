package element

import (
	"html/template"
	"strings"

	"github.com/phper-go/frame/func/conv"
)

func Text(input_name string, input_desc string, input_tip string) Interface {

	elem := &inputElement{
		name:  input_name,
		desc:  input_desc,
		itype: "text",
		tip:   input_tip,
	}
	elem.options = make(map[string]string)
	elem.options["name"] = elem.name
	elem.options["type"] = elem.itype
	return elem
}

func Hidden(input_name string) Interface {
	elem := &inputElement{
		name:  input_name,
		itype: "hidden",
	}
	elem.options = make(map[string]string)
	elem.options["name"] = elem.name
	elem.options["type"] = elem.itype
	return elem
}
func Password(input_name string, input_desc string, input_tip string) Interface {
	elem := &inputElement{
		name:  input_name,
		desc:  input_desc,
		itype: "password",
		tip:   input_tip,
	}
	elem.options = make(map[string]string)
	elem.options["name"] = elem.name
	elem.options["type"] = elem.itype
	return elem
}
func File(input_name string, input_desc string, input_tip string) Interface {
	elem := &inputElement{
		name:  input_name,
		desc:  input_desc,
		itype: "file",
		tip:   input_tip,
	}
	elem.options = make(map[string]string)
	elem.options["name"] = elem.name
	elem.options["type"] = elem.itype
	return elem
}

type inputElement struct {
	name     string
	value    string
	desc     string
	itype    string
	tip      string
	required bool
	options  map[string]string
	errno    string
	errmsg   string
}

func (this *inputElement) Set(value interface{}, required bool) {
	this.value = conv.String(value)
	this.options["value"] = conv.String(this.value)
	this.required = required
}
func (this *inputElement) SetError(errno string, errmsg string) {
	this.errno = errno
	this.errmsg = errmsg
}
func (this *inputElement) Errno() string {
	return this.errno
}
func (this *inputElement) Errmsg() string {
	return this.errmsg
}
func (this *inputElement) HasError() bool {
	if this.errno != "" || this.errmsg != "" {
		return true
	}
	return false
}
func (this *inputElement) Name() string {
	return this.name
}
func (this *inputElement) Desc() string {
	return this.desc
}
func (this *inputElement) Type() string {
	return this.itype
}
func (this *inputElement) Tip() string {
	return this.tip
}
func (this *inputElement) Required() bool {
	return this.required
}
func (this *inputElement) Render(attributes ...string) template.HTML {

	content := "<input "
	for k, v := range this.options {
		content += k + "=\"" + v + "\" "
	}
	for _, attr := range attributes {
		kv := strings.Split(attr, ":")
		if len(kv) == 2 {
			content += kv[0] + "=\"" + kv[1] + "\" "
		} else {
			content += attr
		}
	}
	content += "/>"
	return template.HTML(content)
}
