package route

import (
	"errors"
	"fmt"
	"net/url"
	"strings"

	"github.com/phper-go/frame/func/object"
	"github.com/phper-go/frame/interfaces"
)

var Commands map[string]interface{}
var Controllers map[string]interface{}
var DefaultAction *string
var Debug *int

type Handler struct {
}

func (this *Handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	var http = &HTTP{}
	http.ServeHTTP(response, request)
}

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
	path_len := len(path)
	switch path_len {
	case 0:
		controller = ""
		action = ""
	case 1:
		controller = "/" + path[0]
		action = "index"
	default:
		controller = "/" + strings.Join(path[0:path_len-1], "/")
		action = path[path_len-1]
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
		fmt.Println("not found")
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
