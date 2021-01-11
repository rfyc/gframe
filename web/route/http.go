package route

import (
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"time"

	"github.com/phper-go/frame/web/ctx"

	"github.com/phper-go/frame/func/conv"
	"github.com/phper-go/frame/func/ip"
	"github.com/phper-go/frame/func/object"
	"github.com/phper-go/frame/interfaces"
	"github.com/phper-go/frame/logger"
	"github.com/phper-go/frame/web/session"
)

type httpRouter struct {
	request        *http.Request
	response       http.ResponseWriter
	btime          time.Time
	execController interfaces.Controller
	execAction     string
}

func (this *httpRouter) ServeHTTP(response http.ResponseWriter, request *http.Request) {

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

func (this *httpRouter) paserController() error {

	var err error
	var uri = this.request.RequestURI
	this.execController, this.execAction, err = parseController(uri)
	return err
}

func (this *httpRouter) runController() {

	if this.execAction == "Run" {
		object.Set(this.execController, this.execController.Ctx().Input.Request)
	}

	if isok := this.execController.Prepare(); isok {
		run := reflect.ValueOf(this.execController).MethodByName(this.execAction)
		run.Call([]reflect.Value{})
	}
}

func (this *httpRouter) outputController() {

	//******** set session ********//
	if err := sessionWrite(this.execController.Ctx().Session); err != nil {
		panic(errors.New("write session fail:"))
	}

	//******** set status ********//
	var output = this.execController.Ctx().Output
	var location = conv.String(output.Headers["Location"])
	if len(location) > 0 {
		http.Redirect(this.response, this.request, location, http.StatusFound)
	} else {
		this.write(http.StatusOK, output.Status, output.Error, output.Content)
	}
}

func (this *httpRouter) endController() {

	var output = this.execController.Ctx().Output
	var response = this.response

	//******** set header ********//
	for key, val := range output.Headers {
		response.Header().Set(key, val)
	}
	if this.execController.Ctx().Session.SID == "" {
		this.execController.Ctx().Session.SID = session.ID()
		output.Cookies = append(output.Cookies, &http.Cookie{
			Name:    session.Name,
			Value:   this.execController.Ctx().Session.SID,
			Expires: time.Unix(time.Now().Unix()+int64(session.LifeTime), 0),
		})
	}
	//******** set cookies ********//
	for _, cookie := range output.Cookies {
		http.SetCookie(response, cookie)
	}

	this.execController.End()

}

func (this *httpRouter) initController() error {

	//******** request get ********//

	var Ctx = ctx.NewCtxHTTP(this.response, this.request)
	var _ctx = this.execController.Ctx()
	_ctx.Action = Ctx.Action
	_ctx.Controller = Ctx.Controller
	_ctx.Context = Ctx.Context
	_ctx.Input = Ctx.Input
	_ctx.Output = Ctx.Output
	_ctx.Session = Ctx.Session
	_ctx.Theme = Ctx.Theme
	return sessionRead(_ctx)

}

func (this *httpRouter) error(status int, err error) {

	if *Debug >= 1 {
		this.write(status, "-", err.Error(), []byte(err.Error()))
	} else {
		this.write(status, "-", err.Error(), []byte{})
	}
}

func (this *httpRouter) write(status int, errno, errmsg string, content []byte) {

	this.response.WriteHeader(status)
	size, err := this.response.Write(content)
	this.logger(status, errno, errmsg, size, err)
}

func (this *httpRouter) logger(status int, errno, errmsg string, output_size int, output_err error) {

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
