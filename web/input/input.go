package input

import (
	"net"
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
