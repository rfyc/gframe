package connect

import (
	"errors"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/phper-go/frame/func/file"

	"github.com/phper-go/frame/func/conv"
)

type HTTPServer struct {
	Address  string
	Static   string
	listener net.Listener
	server   *http.Server
	serveMux *http.ServeMux
}

func (this *HTTPServer) Construct() {

	this.serveMux = http.NewServeMux()
	this.server = &http.Server{Handler: this.serveMux}
}

func (this *HTTPServer) Server() *http.Server {
	return this.server
}

func (this *HTTPServer) Listener() net.Listener {
	return this.listener
}

func (this *HTTPServer) listen(addr string) (net.Listener, error) {

	graceful := os.Getenv("graceful")
	if graceful == "on" {
		fd_ptr := conv.Int(os.Getenv(addr))
		fd := os.NewFile(uintptr(fd_ptr), "")
		return net.FileListener(fd)
	}
	return net.Listen("tcp", addr)
}

func (this *HTTPServer) Init() error {

	if this.Address == "" {
		return nil
	}
	listener, err := this.listen(this.Address)
	if err != nil {
		return err
	}
	this.listener = listener
	return nil
}

func (this *HTTPServer) Handle(pattern string, handler http.Handler) {
	this.serveMux.Handle(pattern, handler)
}

func (this *HTTPServer) Start() error {

	if this.listener != nil {

		if this.Static != "" && file.IsDir(this.Static) {
			var basepath = "/" + filepath.Base(this.Static) + "/"
			this.serveMux.Handle(basepath, http.StripPrefix(basepath, http.FileServer(http.Dir(this.Static))))
		}

		return this.server.Serve(this.listener)
	}
	return errors.New("listener empty")
}

func (this *HTTPServer) Destruct() error {
	return nil
}

func (this *HTTPServer) Stop() {

}
