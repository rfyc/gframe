package element

import (
	"html/template"

	"github.com/phper-go/frame/func/conv"
)

func Checkbox(input_name string, input_desc string, items map[string]string, input_tip ...string) Interface {

	elem := &checkboxElement{}
	elem.name = input_name
	elem.desc = input_desc
	elem.itype = "checkbox"
	if len(input_tip) != 0 {
		elem.tip = input_tip[0]
	}
	for k, v := range items {
		e := &inputElement{
			name:  input_name,
			itype: "checkbox",
			value: k,
			desc:  v,
		}
		e.options = make(map[string]string)
		e.options["name"] = e.name
		e.options["type"] = e.itype
		e.options["value"] = e.value
		elem.elems = append(elem.elems, e)
	}
	elem.options = make(map[string]string)
	return elem
}

type checkboxElement struct {
	inputElement
	elems []*inputElement
}

func (this *checkboxElement) Set(value interface{}, required bool) {

	array, ok := value.([]interface{})
	if ok {
		for _, elem := range this.elems {
			for _, s := range array {
				if elem.value == conv.String(s) {
					elem.options["checked"] = "1"
				}
			}
		}
	}
	this.required = required
}
func (this *checkboxElement) Render(attributes ...string) template.HTML {
	content := ""
	for _, elem := range this.elems {
		content += "<label class='checkbox'><input "
		for k, v := range elem.options {
			content += k + "=\"" + v + "\" "
		}
		content += "/> " + elem.desc
		content += " </label>"
	}
	return template.HTML(content)
}
