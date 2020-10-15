package commands

import (
	"github.com/phper-go/frame/web/core"
)

func init() {

	if _, ok := core.ParseCommand("help"); !ok {
		core.RegisterCommand("help", &HelpCmd{})
	}
	if _, ok := core.ParseCommand("start"); !ok {
		core.RegisterCommand("start", &StartCmd{})
	}
}
