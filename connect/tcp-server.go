package connect

import (
	"errors"
	"fmt"
	"net"

	"github.com/phper-go/frame/interfaces"
)

type TCPServer struct {
	Address  string
	listener net.Listener
	handler  interfaces.TCPHandler
}

func (this *TCPServer) Construct() {
}

func (this *TCPServer) Handle(pattern string, handler interfaces.TCPHandler) {
	this.handler = handler
}

func (this *TCPServer) Init() error {
	return nil
}

func (this *TCPServer) Start() error {

	if this.Address == "" {
		return errors.New("address empty")
	}

	ln, err := net.Listen("tcp", this.Address)
	if err != nil {
		return fmt.Errorf("tcp listen error: %v", err)
	}

	this.listener = ln

	for {
		Conn, err := ln.Accept()
		if err == nil {
			fmt.Errorf("tcp accept error: %v", err)
		}
		go func(conn net.Conn) {
			if this.handler != nil {
				this.handler.ServeTCP(conn)
			}

		}(Conn)
	}
	return nil
}

func (this *TCPServer) Close() {
	this.listener.Close()
}
