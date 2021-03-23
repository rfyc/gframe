package element

import (
	"html/template"
	"strings"
)

func Textarea(input_name string, input_desc string, editor bool, input_tip ...string) Interface {

	elem := &textareaElement{}
	elem.name = input_name
	elem.desc = input_desc
	elem.itype = "textarea"
	if len(input_tip) != 0 {
		elem.tip = input_tip[0]
	}
	elem.options = make(map[string]string)
	elem.options["name"] = elem.name
	if editor {
		elem.options["class"] = "editor"
	}
	return elem
}

type textareaElement struct {
	inputElement
}

func (this *textareaElement) Render(attributes ...string) template.HTML {
	content := "<textarea "
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
	content += ">" + this.value
	content += "</textarea>"
	return template.HTML(content)
}
