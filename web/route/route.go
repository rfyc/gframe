package route

import (
	"net/http"
)

var Commands map[string]interface{}
var Controllers map[string]interface{}
var DefaultAction *string
var Debug *int

type Handler struct {
}

func (this *Handler) ServeHTTP(response http.ResponseWriter, request *http.Request) {

	var router = &httpRouter{}
	router.ServeHTTP(response, request)
}

func (this *Handler) ServeTCP() {

}

func init() {
	var debug = 0
	var defaultAction = ""
	Debug = &debug
	DefaultAction = &defaultAction
	Commands = make(map[string]interface{})
	Controllers = make(map[string]interface{})
}
