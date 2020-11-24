package input

import (
	"net"

	"github.com/phper-go/frame/func/object"
)

type Input struct {
	Controller string
	Action     string
	Request    map[string]interface{}
	Get        map[string]interface{}
	Post       map[string]interface{}
	Cookie     map[string]interface{}
	Session    map[string]interface{}
	Header     map[string]interface{}
	Server     *Server
}

func (this *Input) Bind(obj interface{}) error {

	return object.Set(obj, this.Request)
}

type Server struct {
	IsGet         bool
	IsPost        bool
	IsAjax        bool
	IsHttp        bool
	IsHttps       bool
	IsTcp         bool
	IsUdp         bool
	TCPConn       net.TCPConn
	RemoteAddr    string
	ServerName    string
	ServerPort    string
	QueryPath     string
	QueryString   string
	HttpReferer   string
	HttpUserAgent string
}
