package core

import (
	"github.com/phper-go/frame/web/input"
	"github.com/phper-go/frame/web/output"
	"github.com/phper-go/frame/web/session"
	"github.com/phper-go/frame/web/theme"
)

type Controller struct {
	input   *input.Input
	output  *output.Output
	session *session.Session
	theme   theme.Interface
}

func (this *Controller) Theme() theme.Interface {
	return this.theme
}

func (this *Controller) Session() *session.Session {
	return this.session
}

func (this *Controller) Input() *input.Input {
	return this.input
}

func (this *Controller) Output() *output.Output {
	return this.output
}

func (this *Controller) Construct(controllerName string, actionName string) {

	this.input = &input.Input{
		Controller: controllerName,
		Action:     actionName,
		Request:    make(map[string]interface{}),
		Get:        make(map[string]interface{}),
		Post:       make(map[string]interface{}),
		Cookie:     make(map[string]interface{}),
		Session:    make(map[string]interface{}),
		Header:     make(map[string]interface{}),
		Server:     &input.Server{},
	}

	this.output = &output.Output{
		Headers: make(map[string]string),
		Status:  "200",
	}

	this.session = &session.Session{}

	this.theme = &theme.Theme{}
	this.theme.Construct()
}

func (this *Controller) Prepare() bool {
	return true
}

func (this *Controller) End() {

}
