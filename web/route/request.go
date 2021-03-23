package route

import (
	"net/http"
	"strings"
)

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
