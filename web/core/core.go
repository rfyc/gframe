package core

import (
	"strings"

	"github.com/phper-go/frame/interfaces"
	"github.com/phper-go/frame/web/route"
)

var Config interfaces.Config

func RegisterCommand(cmdName string, class interface{}, nx ...bool) {

	if len(nx) > 1 && nx[0] {
		if _, ok := route.Commands[strings.ToLower(cmdName)]; !ok {
			route.Commands[strings.ToLower(cmdName)] = class
		}
	} else {
		route.Commands[strings.ToLower(cmdName)] = class
	}
}

func RegisterController(path string, class interface{}, nx ...bool) {

	if len(nx) > 1 && nx[0] {
		if _, ok := route.Controllers[strings.ToLower(path)]; !ok {
			route.Controllers[strings.ToLower(path)] = class
		}
	} else {
		route.Controllers[strings.ToLower(path)] = class
	}
}

func init() {
	Config = &config{}
}
