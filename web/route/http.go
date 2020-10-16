package route

import (
	"net/http"
	"time"

	"github.com/phper-go/frame/func/conv"

	"github.com/phper-go/frame/logger"

	"github.com/phper-go/frame/interfaces"
)

type HTTP struct {
	DefaultAction  *string
	Debug          *uint8
	request        *http.Request
	response       http.ResponseWriter
	btime          time.Time
	execController interfaces.Controller
	execAction     string
}

func (this *HTTP) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	this.request = request
	this.response = response

	if err := this.paserController(); err != nil {
		this.write(http.StatusNotFound, "-", err.Error(), []byte{})
		return
	}

	if err := this.initController(); err != nil {
		this.write(http.StatusBadGateway, "-", err.Error(), []byte{})
		return
	}

	this.runController()

	this.endController()

	this.outputController()

	defer func() {
		if err := recover(); err != nil {
			// logger.Format(err).Error("run")
			status := http.StatusInternalServerError
			this.write(status, "-", conv.String(err), []byte{})
			return
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

	if isok := this.execController.Prepare(); isok {
		run := reflect.ValueOf(execController).MethodByName(execMethod)
		run.Call([]reflect.Value{})
	}
}

func (this *HTTP) outputController() {

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

func (this *HTTP) endController() {

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

	this.execController.End()

}

func (this *HTTP) initController() {

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

func (this *HTTP) write(status int, errno, errmsg, content []string) {

	this.response.WriteHeader()
	size, err := response.Write(content)
	if err != nil {
		logger.Format(err).Error("run")
	}
	this.logger()
}

func (this *HTTP) debugError(err error) []byte {

	if *this.Debug > 0 {
		return []byte(err.Error())
	}
	return []byte{}
}

func (this *HTTP) logger() {

	logmsg := "http^"
	logmsg += ip.ClientIP(request) + "^"
	logmsg += request.Method + "^"
	logmsg += request.URL.String() + "^"
	// logmsg += request.UserAgent() + "^"
	logmsg += conv.String(outputSize) + "^"
	logmsg += fmt.Sprintf("%0.4f", time.Since(this.beginTime).Seconds()) + "^"
	logmsg += conv.String(status) + "^"
	logmsg += errno

	logFmt := Format(logmsg)
	logFmt.caller = false
	logFmt.Writeln("access")
}
