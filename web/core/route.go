package core

import (
	"encoding/json"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"reflect"
	"strings"
	"time"

	"github.com/phper-go/frame/connect"
	"github.com/phper-go/frame/func/conv"
	"github.com/phper-go/frame/func/object"
	"github.com/phper-go/frame/interfaces"
	"github.com/phper-go/frame/logger"
	"github.com/phper-go/frame/web/input"
	"github.com/phper-go/frame/web/output"
	"github.com/phper-go/frame/web/session"
)

type route struct {
	DefaultAction *string
	Debug         *uint8
}

func (this *route) findController(uri string) (execController interfaces.Controller, execMethod string, err error) {

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

func (this *route) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	var access = logger.Access()
	defer func() {
		if err := recover(); err != nil {
			logger.Format(err).Error("run")
			response.WriteHeader(http.StatusInternalServerError)
			if *this.Debug > 0 {
				this.writeHTTP(response, conv.Bytes(err))
			}
			access.HTTP(request, response, http.StatusInternalServerError, 0, "-")
			return
		}
		return
	}()

	//******** parse uri ********//
	execController, execMethod, err := this.findController(request.RequestURI)
	if err != nil {
		response.WriteHeader(http.StatusNotFound)
		size := this.writeHTTP(response, conv.Bytes(err.Error()))
		access.HTTP(request, response, http.StatusNotFound, size, "-")
		return
	}

	//******** controller init ********//
	input := execController.Input()
	output := execController.Output()
	this.initHTTP(input, request, response)
	object.Set(execController, input.Request)

	//******** session read ********//
	if err := this.readSession(input, output); err != nil {
		response.WriteHeader(http.StatusBadGateway)
		content := conv.Bytes("get session bad")
		size := this.writeHTTP(response, content)
		access.HTTP(request, response, http.StatusBadGateway, size, "get_session_bad")
		return
	}

	//******** action run ********//
	if isok := execController.Prepare(); isok {
		run := reflect.ValueOf(execController).MethodByName(execMethod)
		run.Call([]reflect.Value{})
	}

	execController.End()

	//******** set header ********//
	for key, val := range output.Headers {
		response.Header().Set(key, val)
	}
	//******** set cookies ********//
	for _, cookie := range output.Cookies {
		http.SetCookie(response, cookie)
	}
	//******** set session ********//
	if err := this.writeSession(input); err != nil {
		access.HTTP(request, response, http.StatusBadGateway, 0, "set_session_bad")
	}

	//******** set status ********//
	var size = 0
	var location = conv.String(output.Headers["Location"])
	if output.Status == http.StatusFound || output.Status == http.StatusSeeOther {
		http.Redirect(response, request, location, output.Status)
	} else {
		response.WriteHeader(output.Status)
		size = this.writeHTTP(response, output.Content)
	}
	//******** set output ********//
	access.HTTP(request, response, output.Status, size, output.Error)
}

func (this *route) ServeTCP(conn net.Conn) {

	var session = make(map[string]interface{})

	for {
		var access = logger.Access()
		var content = &connect.TCPContent{}
		var body = make(map[string]interface{})
		var output_content = &connect.TCPContent{}
		output_content.Type = connect.TCP_PROTOCOL_V1
		output_content.Protocol = connect.TCP_TYPE_RESPONSE
		//******** parse package ********//
		if _, err := content.Recv(conn); err != nil {
			if err == io.EOF {
				conn.Close()
				access.TCP(conn, body, 0, "510", "close")
				return
			}
			output_content.Body, _ = json.Marshal(map[string]string{
				"errno":  "511",
				"errmsg": "recv error: " + err.Error(),
			})
			size, send_err := output_content.Send(conn)
			if send_err != nil {
				err = errors.New(err.Error() + " - send:" + send_err.Error())
			}
			access.TCP(conn, body, size, "511", "recv error: "+err.Error())
			continue
		}

		//******** parse uri ********//
		if err := json.Unmarshal(content.Body, &body); err != nil {
			output_content.Body, _ = json.Marshal(map[string]string{
				"errno":  "512",
				"errmsg": "json decode error: " + err.Error(),
			})
			size, send_err := output_content.Send(conn)
			if send_err != nil {
				err = errors.New(err.Error() + " - send:" + send_err.Error())
			}
			access.TCP(conn, body, size, "512", "json decode error: "+err.Error())
			continue
		}
		uri := conv.String(body["_api_"])
		controllerName, actionName := this.parseAction(uri)

		//******** found controller ********//
		execController, err := ParseController(controllerName)
		if err != nil {
			output_content.Body, _ = json.Marshal(map[string]string{
				"errno":  "513",
				"errmsg": "parse controller: " + err.Error(),
			})
			size, send_err := output_content.Send(conn)
			if send_err != nil {
				err = errors.New(err.Error() + " - send:" + send_err.Error())
			}
			access.TCP(conn, body, size, "513", "parse controller: "+err.Error())
			continue
		}

		//******** found action ********//
		execMethod := object.FindMethod(execController, actionName+"Action")
		if execMethod == "" {
			output_content.Body, _ = json.Marshal(map[string]string{
				"errno":  "514",
				"errmsg": "parse action: " + err.Error(),
			})
			size, send_err := output_content.Send(conn)
			if send_err != nil {
				err = errors.New(err.Error() + " - send:" + send_err.Error())
			}
			access.TCP(conn, body, size, "513", "parse action: "+err.Error())
			continue
		}
		//******** controller init ********//
		execController.Construct(controllerName, actionName)
		input := execController.Input()
		input.Get = body
		input.Post = body
		input.Request = body
		input.Server.IsTcp = true
		input.Server.RemoteAddr = conn.RemoteAddr().String()
		input.Session = session

		//******** action run ********//
		if isok := execController.Prepare(); isok {
			run := reflect.ValueOf(execController).MethodByName(execMethod)
			run.Call([]reflect.Value{})
		}

		execController.End()

		output := execController.Output()
		output_content = &connect.TCPContent{}
		output_content.Type = connect.TCP_PROTOCOL_V1
		output_content.Protocol = connect.TCP_TYPE_RESPONSE
		output_content.Body = output.Content
		size, err := output_content.Send(conn)
		if err != nil {
			output.Error += " - send:" + err.Error()
		}
		access.TCP(conn, body, size, "200", output.Error)

	}

}

func (this *route) readSession(input *input.Input, output *output.Output) error {

	if session.Enable > 0 {
		session_id := conv.String(input.Cookie[session.Name])
		if len(session_id) == 0 {
			session_id = session.ID()
			output.Cookies = append(output.Cookies, &http.Cookie{
				Name:    session.Name,
				Value:   session_id,
				Expires: time.Unix(time.Now().Unix()+int64(session.LifeTime), 0),
			})
			input.Cookie[session.Name] = session_id
			return nil
		}
		sessioin_data, err := session.Read(session_id)
		if err != nil {
			return errors.New("session read error: " + err.Error())
		}
		input.Session = sessioin_data
	}
	return nil
}

func (this *route) writeSession(input *input.Input) error {
	session_id := conv.String(input.Cookie[session.Name])
	return session.Write(session_id, input.Session)
}

func (this *route) initHTTP(input *input.Input, request *http.Request, response http.ResponseWriter) {

	//******** request get ********//
	for key, value := range request.URL.Query() {
		count := len(value)
		if count > 1 {
			input.Request[key] = value
			input.Get[key] = value
		} else if count == 1 {
			input.Request[key] = value[0]
			input.Get[key] = value[0]
		}
	}

	request.ParseForm()

	//******** request post ********//
	for key, value := range request.PostForm {
		count := len(value)
		if count > 1 {
			input.Request[key] = value
			input.Post[key] = value
		} else if count == 1 {
			input.Request[key] = value[0]
			input.Post[key] = value[0]
		}
	}

	//******** header ********//
	for key, value := range request.Header {
		count := len(value)
		if count > 1 {
			input.Header[key] = value
		} else if count == 1 {
			input.Header[key] = value[0]
		}
	}

	//******** cookie ********//
	for _, cookie := range request.Cookies() {
		input.Cookie[cookie.Name] = cookie.Value
	}

	//******** server ********//
	input.Server.IsHttp = true
	input.Server.IsPost = request.Method == "POST"
	input.Server.IsGet = request.Method == "GET"
	input.Server.IsAjax = len(request.Header.Get("x-requested-with")) > 0
	input.Server.RemoteAddr = request.RemoteAddr
	input.Server.ServerName = this.httpServerName(request)
	input.Server.ServerPort = this.httpServerPort(request)
	input.Server.QueryPath = this.httpQueryPath(request)
	input.Server.QueryString = this.httpQueryString(request)
	input.Server.HttpReferer = request.Referer()
	input.Server.HttpUserAgent = request.UserAgent()

}

func (this *route) httpQueryString(request *http.Request) string {
	params := strings.Split(request.RequestURI, "?")
	if len(params) > 1 {
		return params[1]
	}
	return ""
}

func (this *route) httpQueryPath(request *http.Request) string {
	params := strings.Split(request.RequestURI, "?")
	if len(params) > 0 {
		return params[0]
	}
	return ""
}

func (this *route) httpServerName(request *http.Request) string {

	var host string
	if len(request.Host) > 0 {
		host = request.Host
	} else if len(request.URL.Host) > 0 {
		host = request.URL.Host
	}
	params := strings.Split(host, ":")
	if len(params) > 0 {
		return params[0]
	}
	return ""
}

func (this *route) httpServerPort(request *http.Request) string {

	var host string
	if len(request.Host) > 0 {
		host = request.Host
	} else if len(request.URL.Host) > 0 {
		host = request.URL.Host
	}
	params := strings.Split(host, ":")
	if len(params) > 1 {
		return params[1]
	}
	return ""
}

func (this *route) writeHTTP(response http.ResponseWriter, content []byte) (size int) {

	size, err := response.Write(content)
	if err != nil {
		logger.Format(err).Error("run")
	}
	return size
}

func (this *route) parseURI(requestURI string) (controller, action string, err error) {

	urlPath, err := url.ParseRequestURI(requestURI)
	if err != nil {
		return "", "", err
	}

	controller, action = this.parseAction(urlPath.Path)
	return controller, action, nil
}

func (this *route) parseAction(uri string) (controller, action string) {

	uri = strings.Trim(uri, "/")
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
	return strings.ToLower(controller), strings.ToLower(action)
}
