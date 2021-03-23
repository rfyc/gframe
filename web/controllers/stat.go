package controllers

import (
	"os"

	"github.com/phper-go/frame/web/core"
)

type StatController struct {
	core.Controller
}

func (this *StatController) Prepare() bool {

	if ok := PrepareHandle(this); !ok {
		return ok
	}
	return true
}

func (this *StatController) InfoAction() {

	this.Output.Echo(os.Getpid(), "-", "-")
}
