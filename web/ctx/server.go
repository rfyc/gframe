package ctx

type Server struct {
	IsGet         bool
	IsPost        bool
	IsAjax        bool
	IsHttp        bool
	IsHttps       bool
	IsTcp         bool
	IsUdp         bool
	ServerName    string
	ServerPort    string
	QueryPath     string
	QueryString   string
	HttpReferer   string
	HttpUserAgent string
	RemoteAddr    string
}
