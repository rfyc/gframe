package route

func parseURI(requestURI string) (controller, action string, err error) {

	urlPath, err := url.ParseRequestURI(requestURI)
	if err != nil {
		return "", "", err
	}

	uri := strings.Trim(urlPath.Path, "/")
	if len(uri) == 0 {
		uri = strings.Trim(*this.DefaultAction, "/")
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
	controllerName, actionName, err := this.parseURI(uri)
	if err != nil {
		return nil, "", errors.New(uri + " parse fail")
	}

	//******** found controller ********//
	execController, err = ParseController(controllerName + "/" + actionName)
	if err == nil {
		if _, ok := execController.(interfaces.Action); ok {
			execController.Construct(controllerName, actionName)
			return execController, "Run", nil
		}
	}

	//******** found controller ********//
	execController, err = ParseController(controllerName)
	if err != nil {
		return nil, "", errors.New(uri + " not found")
	}

	//******** found action ********//
	execMethod = object.FindMethod(execController, actionName+"Action")
	if execMethod == "" {
		return nil, "", errors.New(uri + " action not found")
	}

	//******** controller init ********//
	execController.Construct(controllerName, actionName)

	return execController, execMethod, nil
}
