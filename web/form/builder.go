package form

import (
	"reflect"
	"strings"

	"github.com/phper-go/frame/ext/validator"
	"github.com/phper-go/frame/func/object"
	"github.com/phper-go/frame/web/form/element"
)

func BuildForm(action, method string, elements []element.Interface, formModel ...interface{}) *Builder {

	builder := &Builder{}

	builder.Construct(action, method)

	builder.Elements = elements

	if len(formModel) > 0 {

		builder.Bind(formModel[0])

	}

	builder.Enctype = "multipart/form-data"

	return builder
}

type Builder struct {
	Title    string
	Method   string
	Action   string
	Target   string
	Enctype  string
	Elements []element.Interface
}

func (this *Builder) Construct(action string, method string) {
	this.Action = action
	this.Method = method
	this.Elements = []element.Interface{}
}

func (this *Builder) Bind(obj interface{}) {

	mapRequireds := make(map[string]bool)
	valiObj, is_vali := obj.(validator.Interface)
	if is_vali {
		rules := valiObj.Rules()
		for _, rule := range rules {
			if reflect.TypeOf(rule).String() == reflect.TypeOf(&validator.Required{}).String() {
				fields := strings.Split(rule.GetFields(), ",")
				for _, field := range fields {
					mapRequireds[strings.ToLower(field)] = true
					continue
				}
			}
		}
	}

	var errno, errmsg, field string
	var values = object.Values(obj)
	if tmp, ok := obj.(validator.ApiInterface); ok {
		errno, errmsg, field = tmp.GetErrors()
	}
	for _, elem := range this.Elements {
		var fieldVal interface{}
		elemName := strings.ToLower(elem.Name())
		for key, val := range values {
			if strings.ToLower(key) == elemName {
				fieldVal = val
			}
		}
		if !(strings.ToLower(elem.Type()) == "submit" || strings.ToLower(elem.Type()) == "button") {
			if strings.ToLower(elem.Type()) == "widget" {
				fieldVal = obj
			}
			elem.Set(fieldVal, mapRequireds[elemName])
		}

		if field != "" && strings.ToLower(field) == elemName {
			elem.SetError(errno, errmsg)
		}
	}
}
