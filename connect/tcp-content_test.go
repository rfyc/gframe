package connect

import (
	"net"
	"testing"
	"time"
)

func Test_Content(t *testing.T) {

	ln, err := net.Listen("tcp", ":8090")
	if err != nil {
		t.Error("tcp listen error: ", err.Error())
		return
	}

	go func() {
		for {

			Conn, err := ln.Accept()
			if err != nil {
				t.Error("tcp accept error: ", err.Error())
			}

			t.Log(Conn.RemoteAddr())

			go func(conn net.Conn) {
				for {
					content := &TCPContent{}
					if _, err := content.Recv(conn); err != nil {
						t.Log("tcp recv error:", err.Error())
						return
					}
					t.Log("tcp recv protocol:", content.Protocol)
					t.Log("tcp recv type:", content.Type)
					t.Log("tcp recv:", string(content.Body))
				}
			}(Conn)
		}
	}()

	conn, err := net.Dial("tcp", "127.0.0.1:8090")
	if err != nil {
		t.Error("net error:", err.Error())
		return
	}

	content := &TCPContent{}

	content.Type = TCP_TYPE_REQUEST

	content.Protocol = TCP_PROTOCOL_V1

	content.Body = []byte("testestsedtset")

	t.Log(content.Send(conn))

	time.Sleep(1 * time.Second)
	t.Log(content.Send(conn))

	time.Sleep(5 * time.Second)
	t.Log(content.Send(conn))

	time.Sleep(1 * time.Second)
	t.Log(content.Send(conn))

	conn.Close()
}
