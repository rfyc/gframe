package element

import (
	"html/template"
)

type Interface interface {
	Set(value interface{}, required bool)
	SetError(errno string, errmsg string)
	HasError() bool
	Errno() string
	Errmsg() string
	Name() string
	Desc() string
	Type() string
	Tip() string
	Required() bool
	Render(attributes ...string) template.HTML
}
