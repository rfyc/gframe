package ctx

import (
	"net/http"
	"strings"

	"github.com/phper-go/frame/func/object"
)

type Input struct {
	Request map[string]interface{}
	Get     map[string]interface{}
	Post    map[string]interface{}
	Cookie  map[string]interface{}
	Header  map[string]interface{}
}

func (this *Input) Bind(obj interface{}) error {

	return object.Set(obj, this.Request)
}

func NewInput() *Input {
	return &Input{
		Request: make(map[string]interface{}),
		Get:     make(map[string]interface{}),
		Post:    make(map[string]interface{}),
		Header:  make(map[string]interface{}),
		Cookie:  make(map[string]interface{}),
	}
}

func NewInputHTTP(request *http.Request) *Input {

	var input = NewInput()

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

	//******** request post ********//
	request.ParseForm()
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

	return input
}

func NewServerHTTP(request *http.Request) *Server {

	sever := &Server{}
	sever.IsHttp = true
	sever.IsPost = request.Method == "POST"
	sever.IsGet = request.Method == "GET"
	sever.IsAjax = len(request.Header.Get("x-requested-with")) > 0
	sever.RemoteAddr = request.RemoteAddr
	sever.ServerName = httpServerName(request)
	sever.ServerPort = httpServerPort(request)
	sever.QueryPath = httpQueryPath(request)
	sever.QueryString = httpQueryString(request)
	sever.HttpReferer = request.Referer()
	sever.HttpUserAgent = request.UserAgent()

	return sever
}

func httpQueryString(request *http.Request) string {
	params := strings.Split(request.RequestURI, "?")
	if len(params) > 1 {
		return params[1]
	}
	return ""
}

func httpQueryPath(request *http.Request) string {
	params := strings.Split(request.RequestURI, "?")
	if len(params) > 0 {
		return params[0]
	}
	return ""
}

func httpServerName(request *http.Request) string {

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

func httpServerPort(request *http.Request) string {

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
