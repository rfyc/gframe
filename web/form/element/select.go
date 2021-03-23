package element

import (
	"html/template"
)

func Select(input_name string, input_desc string, option_value_descs map[string]string, input_tip ...string) Interface {
	elem := &selectElement{}
	elem.name = input_name
	elem.desc = input_desc
	elem.itype = "select"
	if len(input_tip) != 0 {
		elem.tip = input_tip[0]
	}
	for k, v := range option_value_descs {
		e := &inputElement{
			name:  input_name,
			itype: "option",
			value: k,
			desc:  v,
		}
		e.options = make(map[string]string)
		e.options["value"] = e.value
		elem.elems = append(elem.elems, e)
	}
	elem.options = make(map[string]string)
	return elem
}

type selectElement struct {
	inputElement
	elems []*inputElement
}

func (this *selectElement) Set(value interface{}, required bool) {
	s, ok := value.(string)
	if ok {
		for _, elem := range this.elems {
			if elem.value == s {
				elem.options["selected"] = "1"
			}
		}
	}
	this.required = required
}
func (this *selectElement) Render(attributes ...string) template.HTML {

	content := "<select name=\"" + this.name + "\" >"
	for _, elem := range this.elems {
		content += "<option "
		for k, v := range elem.options {
			content += k + "=\"" + v + "\" "
		}
		content += "> " + elem.desc
		content += " </option>"
	}
	content += "</select>"
	return template.HTML(content)
}
