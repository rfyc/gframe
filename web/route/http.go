package route

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/phper-go/frame/func/conv"
	"github.com/phper-go/frame/func/ip"
	"github.com/phper-go/frame/func/object"
	"github.com/phper-go/frame/interfaces"
	"github.com/phper-go/frame/logger"
	"github.com/phper-go/frame/web/session"
)

type HTTP struct {
	request        *http.Request
	response       http.ResponseWriter
	btime          time.Time
	execController interfaces.Controller
	execAction     string
}

func (this *HTTP) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	this.btime = time.Now()
	this.request = request
	this.response = response

	if err := this.paserController(); err != nil {
		this.error(http.StatusNotFound, err)
		return
	}

	if err := this.initController(); err != nil {
		this.error(http.StatusBadGateway, err)
		return
	}

	this.runController()

	this.endController()

	this.outputController()

	defer func() {
		if err := recover(); err != nil {
			// logger.Format(err).Error("run")
			this.error(http.StatusInternalServerError, errors.New(conv.String(err)))
		}
		return
	}()

}

func (this *HTTP) paserController() error {

	var err error
	var uri = this.request.RequestURI
	this.execController, this.execAction, err = parseController(uri)
	return err
}

func (this *HTTP) runController() {

	if this.execAction == "Run" {
		object.Set(this.execController, this.execController.Input().Request)
	}

	if isok := this.execController.Prepare(); isok {
		run := reflect.ValueOf(this.execController).MethodByName(this.execAction)
		run.Call([]reflect.Value{})
	}
}

func (this *HTTP) outputController() {

	//******** set session ********//
	if err := sessionWrite(this.execController.Session()); err != nil {
		panic(errors.New("write session fail:"))
	}

	//******** set status ********//
	var output = this.execController.Output()
	var location = conv.String(output.Headers["Location"])
	if len(location) > 0 {
		http.Redirect(this.response, this.request, location, http.StatusFound)
	} else {
		this.write(http.StatusOK, output.Status, output.Error, output.Content)
	}
}

func (this *HTTP) endController() {

	var output = this.execController.Output()
	var response = this.response

	//******** set header ********//
	for key, val := range output.Headers {
		response.Header().Set(key, val)
	}
	if this.execController.Session().SID == "" {
		this.execController.Session().SID = session.ID()
		output.Cookies = append(output.Cookies, &http.Cookie{
			Name:    session.Name,
			Value:   this.execController.Session().SID,
			Expires: time.Unix(time.Now().Unix()+int64(session.LifeTime), 0),
		})
	}
	//******** set cookies ********//
	for _, cookie := range output.Cookies {
		http.SetCookie(response, cookie)
	}

	this.execController.End()

}

func (this *HTTP) initController() error {

	//******** request get ********//
	var request = this.request
	var input = this.execController.Input()
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
	input.Server.ServerName = httpServerName(request)
	input.Server.ServerPort = httpServerPort(request)
	input.Server.QueryPath = httpQueryPath(request)
	input.Server.QueryString = httpQueryString(request)
	input.Server.HttpReferer = request.Referer()
	input.Server.HttpUserAgent = request.UserAgent()

	sessionRead(this.execController)

	return nil
}

func (this *HTTP) error(status int, err error) {

	if *Debug >= 1 {
		this.write(status, "-", err.Error(), []byte(err.Error()))
	} else {
		this.write(status, "-", err.Error(), []byte{})
	}
}

func (this *HTTP) write(status int, errno, errmsg string, content []byte) {

	this.response.WriteHeader(status)
	size, err := this.response.Write(content)
	this.logger(status, errno, errmsg, size, err)
}

func (this *HTTP) logger(status int, errno, errmsg string, output_size int, output_err error) {

	var request = this.request

	logmsg := this.btime.Format("2006-01-02 15:04:05.0000") + " "
	logmsg += "http "
	logmsg += ip.ClientIP(request) + " "
	logmsg += request.Method + " "
	logmsg += fmt.Sprintf("%0.4f", time.Since(this.btime).Seconds()) + " "
	logmsg += conv.String(status) + " "
	logmsg += conv.String(output_size) + " "
	if output_err == nil {
		logmsg += " - | "
	} else {
		logmsg += conv.String(output_err) + " | "
	}
	logmsg += request.URL.String() + " "
	logmsg += errno + " "
	logmsg += errmsg + " "
	logger.Format(logmsg).Writeln("access")
}
