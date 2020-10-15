package element

import (
	"html/template"
	"strings"
)

func Button(btnType string, desc string) Interface {

	elem := &buttonElement{}
	elem.itype = strings.ToLower(btnType)
	elem.value = desc

	elem.options = make(map[string]string)
	elem.options["type"] = btnType
	return elem
}

type buttonElement struct {
	inputElement
}

func (this *buttonElement) Render(attributes ...string) template.HTML {
	content := "<button "
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
	content += ">" + this.value + "</button>"
	return template.HTML(content)
}
