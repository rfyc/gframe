package core

import (
	"github.com/phper-go/frame/ext/validator"
	"github.com/phper-go/frame/interfaces"
)

type Command struct {
	app interfaces.App
}

func (this *Command) Construct(runApp interfaces.App) {
	this.app = runApp
	this.app.Construct()
}
func (this *Command) App() interfaces.App {
	return this.app
}

func (this *Command) Init() error {

	return nil
}

func (this *Command) LoadApp() error {

	if err := Config.Load(this.App()); err != nil {
		return err
	}
	return nil
}

func (this *Command) InitApp() error {

	if err := this.App().Init(); err != nil {
		return err
	}
	return nil
}

func (this *Command) Prepare() error {

	return nil
}

func (this *Command) Rules() interfaces.ValidatorRules {
	return interfaces.ValidatorRules{
		&validator.LocalFile{Fields: "config,logpath", AllowEmpty: true},
	}
}

func (this *Command) Run() {

}

func (this *Command) End() {

}
