package route

import (
	"errors"
	"net/url"
	"strings"

	"github.com/phper-go/frame/func/object"
	"github.com/phper-go/frame/interfaces"
)

var Commands map[string]interface{}
var Controllers map[string]interface{}
var DefaultAction *string
var Debug *int

func parseURI(requestURI, defaultAction string) (controller, action string, err error) {

	urlPath, err := url.ParseRequestURI(requestURI)
	if err != nil {
		return "", "", err
	}

	uri := strings.Trim(urlPath.Path, "/")
	if len(uri) == 0 {
		uri = strings.Trim(defaultAction, "/")
	}

	path := strings.Split(uri, "/")
	switch len(path) {
	case 3:
		controller = "/" + path[0] + "/" + path[1]
		action = path[2]
	case 2:
		controller = "/" + path[0]
		action = path[1]
	case 1:
		controller = "/" + path[0]
		action = "index"
	}

	return strings.ToLower(controller), strings.ToLower(action), nil
}

func parseController(uri string) (execController interfaces.Controller, execMethod string, err error) {

	//******** parse uri ********//
	controllerName, actionName, err := parseURI(uri, *DefaultAction)
	if err != nil {
		return nil, "", errors.New(uri + " parse fail")
	}

	//******** found action ********//
	var action = controllerName + "/" + actionName
	if obj, ok := Controllers[strings.ToLower(action)]; ok {
		if execController, ok = object.New(obj).(interfaces.Action); ok {
			if _, ok := execController.(interfaces.Action); ok {

				execController.Construct(controllerName, actionName)
				return execController, "Run", nil
			}
		}
	}

	//******** found controller ********//
	var ok bool
	var obj interface{}
	if obj, ok = Controllers[strings.ToLower(controllerName)]; !ok {
		return nil, "", errors.New(uri + " not found")
	}

	if execController, ok = object.New(obj).(interfaces.Controller); !ok {
		return nil, "", errors.New(uri + " not controller")
	}

	execMethod = object.FindMethod(execController, actionName+"Action")
	if execMethod == "" {
		return nil, "", errors.New(uri + " action not found")
	}

	//******** controller init ********//
	execController.Construct(controllerName, actionName)

	return execController, execMethod, nil
}

func init() {
	var debug = 0
	var defaultAction = ""
	Debug = &debug
	DefaultAction = &defaultAction
	Commands = make(map[string]interface{})
	Controllers = make(map[string]interface{})
}
