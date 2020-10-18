package commands

import (
	"github.com/phper-go/frame/web/core"
	"github.com/phper-go/frame/web/route"
)

func init() {

	if _, ok := route.Commands["help"]; !ok {
		core.RegisterCommand("help", &HelpCmd{})
	}
	if _, ok := route.Commands["start"]; !ok {
		core.RegisterCommand("start", &StartCmd{})
	}
}
