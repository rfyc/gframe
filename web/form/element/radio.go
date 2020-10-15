package element

import (
	"html/template"
)

func Radio(input_name string, input_desc string, items map[string]string, input_tip ...string) Interface {

	elem := &radioElement{}
	elem.name = input_name
	elem.desc = input_desc
	elem.itype = "radio"
	if len(input_tip) != 0 {
		elem.tip = input_tip[0]
	}
	for k, v := range items {
		e := &inputElement{
			name:  input_name,
			itype: "radio",
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

type radioElement struct {
	inputElement
	elems []*inputElement
}

func (this *radioElement) Set(value interface{}, required bool) {
	val, ok := value.(string)
	if ok {
		for _, elem := range this.elems {
			if elem.value == val {
				elem.options["checked"] = "1"
			}
		}
	}
	this.required = required
}
func (this *radioElement) Render(attributes ...string) template.HTML {
	content := "<div>"
	for _, elem := range this.elems {
		content += "<label class='radio'><input "
		for k, v := range elem.options {
			content += k + "=\"" + v + "\" "
		}

		content += "/> " + elem.desc
		content += " </label>"
	}
	content += "</div>"
	return template.HTML(content)
}
