package connect

import (
	"bytes"
	"encoding/binary"
	"net"

	"github.com/phper-go/frame/func/conv"
)

const (
	TCP_PROTOCOL_V1   = 1
	TCP_PROTOCOL_V2   = 2
	TCP_TYPE_REQUEST  = 1
	TCP_TYPE_RESPONSE = 2
)

type TCPContent struct {
	Protocol uint8
	Type     uint8
	Body     []byte
}

func (this *TCPContent) Recv(conn net.Conn) (size uint32, err error) {

	var len int
	var pkg_len_buf = make([]byte, 8)
	if _, err = conn.Read(pkg_len_buf); err != nil {
		return
	}

	size = binary.BigEndian.Uint32(pkg_len_buf[0:4])
	this.Protocol = pkg_len_buf[4:5][0]
	this.Type = pkg_len_buf[5:6][0]
	this.Body = make([]byte, size)
	var read_len int
	for {
		content := this.Body[read_len:size]
		len, err = conn.Read(content)
		if err != nil {
			return
		}
		read_len += len
		if read_len >= conv.Int(size) {
			return
		}
	}
}

func (this *TCPContent) Send(conn net.Conn) (size int, err error) {

	var buf bytes.Buffer
	var buf_pref = make([]byte, 8)
	binary.BigEndian.PutUint32(buf_pref, uint32(len(this.Body)))
	buf_pref[4] = this.Protocol
	buf_pref[5] = this.Type
	buf.Write(buf_pref)
	buf.Write(this.Body)

	return conn.Write(buf.Bytes())
}
