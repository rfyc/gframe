package commands

import (
	"fmt"

	"github.com/phper-go/frame/web/core"
)

type HelpCmd struct {
	core.Command
}

func (this *HelpCmd) Run() {
	fmt.Println("HelpCmd Run")
}
