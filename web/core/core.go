package core

import (
	"errors"
	"strings"

	"github.com/phper-go/frame/func/object"
	"github.com/phper-go/frame/interfaces"
)

var AppName string
var Config interfaces.Config
var commands map[string]interface{}
var controllers map[string]interface{}

func RegisterCommand(cmdName string, class interface{}, nx ...bool) {

	if len(nx) > 1 && nx[0] {
		if _, ok := commands[strings.ToLower(cmdName)]; !ok {
			commands[strings.ToLower(cmdName)] = class
		}
	} else {
		commands[strings.ToLower(cmdName)] = class
	}
}

func RegisterController(path string, class interface{}, nx ...bool) {

	if len(nx) > 1 && nx[0] {
		if _, ok := controllers[strings.ToLower(path)]; !ok {
			controllers[strings.ToLower(path)] = class
		}
	} else {
		controllers[strings.ToLower(path)] = class
	}
}

func ParseCommand(cmdName string) (interfaces.Command, bool) {

	if cmdClass, ok := commands[strings.ToLower(cmdName)]; ok {
		if execCommand, ok := object.New(cmdClass).(interfaces.Command); ok {
			return execCommand, true
		}
		return nil, false
	}
	return nil, false
}

func ParseController(controllerName string) (interfaces.Controller, error) {

	if obj, ok := controllers[strings.ToLower(controllerName)]; ok {
		if execController, ok := object.New(obj).(interfaces.Controller); ok {
			return execController, nil
		}
		return nil, errors.New("cotroller not exist")
	}
	return nil, errors.New("cotroller not exist")
}

func init() {
	commands = make(map[string]interface{})
	controllers = make(map[string]interface{})
	Config = &config{}
}
