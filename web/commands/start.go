package commands

import (
	"os"

	"github.com/phper-go/frame/logger"
	"github.com/phper-go/frame/web/core"
)

type StartCmd struct {
	core.Command
}

func (this *StartCmd) Run() {

	if err := this.App().Start(); err != nil {
		logger.Format(err.Error(), this).Echo(os.Stderr)
	}

	this.App().Wait()
}
