package logger

import (
	"fmt"
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/phper-go/frame/func/conv"
	"github.com/phper-go/frame/func/ip"
)

type accessInterface interface {
	HTTP(request *http.Request, response http.ResponseWriter, status int, outputSize int, errno string)
	TCP(conn net.Conn, body map[string]interface{}, size int, errno, errmsg string)
}

type access struct {
	beginTime time.Time
}

func (this *access) HTTP(request *http.Request, response http.ResponseWriter, status int, outputSize int, errno string) {

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

func (this *access) TCP(conn net.Conn, body map[string]interface{}, size int, errno, errmsg string) {

	var query string
	for k, v := range body {
		query += "&" + k + "=" + conv.String(v)
	}
	query = strings.Trim(query, "&")

	logmsg := "tcp^"
	logmsg += conn.RemoteAddr().String() + "^"
	logmsg += query + "^"
	logmsg += conv.String(size) + "^"
	logmsg += fmt.Sprintf("%0.4f", time.Since(this.beginTime).Seconds()) + "^"
	logmsg += errno + "^"
	logmsg += errmsg
	logFmt := Format(logmsg)
	logFmt.caller = false
	logFmt.Writeln("access")
}
