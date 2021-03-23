package connect

// import (
// 	"encoding/json"
// 	"fmt"
// 	"net"

// 	"github.com/webapp-go/frame/logger"
// )

// import (
// 	"bytes"
// 	"encoding/binary"
// 	"errors"
// 	"fmt"
// 	"io"
// 	"net"
// 	"strings"
// 	"sync"
// 	"sync/atomic"
// 	"time"

// 	"github.com/webapp-go/frame/util/format"
// )

// var recv_type_query = "0"
// var recv_type_result = "1"

// type TCPConn struct {
// 	Async       bool
// 	ReadTimeout int
// 	Conn        net.Conn
// 	once        *sync.Once
// 	chanResult  *sync.Map
// 	chanQuery   chan string
// 	reqNo       uint32
// 	reading     bool
// }

// func (this *TCPConn) Construct() {
// 	this.once = &sync.Once{}
// 	this.chanResult = &sync.Map{}
// 	this.chanQuery = make(chan string, 64)
// 	this.ReadTimeout = 3
// }

// func (this *TCPConn) recv() {

// 	for {

// 		content, err := Read(this.Conn)
// 		if len(content) < 17 {
// 			continue
// 		}
// 		recv_type := content[0:1]
// 		if string(recv_type) == recv_type_query {
// 			this.chanQuery <- content[1:]
// 		} else if string(recv_type) == recv_type_result {

// 			if chanObj, ok := this.chanResult.Load(content[1:17]); ok {
// 				if readChan, ok := chanObj.(chan string); ok {

// 					readChan <- content[17:]
// 				} else {
// 					fmt.Println("recv readchan nil")
// 				}
// 			} else {
// 				fmt.Println("recv chanobj nil")
// 			}
// 		} else {
// 			fmt.Println("recv_type_nil")
// 		}

// 		if err == io.EOF {
// 			this.chanQuery <- "error:" + err.Error()
// 			break
// 		}
// 		if err != nil {

// 		}

// 	}
// }

// func (this *TCPConn) Request(content []byte) (result string, err error) {

// 	reqNo, _, err := this.SendQuery(content)
// 	if err != nil {
// 		return "", err
// 	}
// 	return this.RecvResult(reqNo)
// }

// func (this *TCPConn) RecvQuery() (reqNo string, content string, err error) {

// 	this.once.Do(func() {
// 		go this.recv()
// 	})

// 	select {
// 	case content = <-this.chanQuery:
// 		if content[:6] == "error:" {
// 			return "", "", errors.New(content[6:])
// 		}
// 		return content[0:16], content[16:], nil
// 	}
// 	return
// }

// func (this *TCPConn) RecvResult(reqNo string) (content string, err error) {

// 	this.once.Do(func() {
// 		go this.recv()
// 	})

// 	if chanObj, ok := this.chanResult.Load(reqNo); ok {
// 		if readChan, ok := chanObj.(chan string); ok {
// 			select {
// 			case content = <-readChan:
// 			case <-time.After(time.Duration(this.ReadTimeout) * time.Second):
// 				err = errors.New("read timeout")
// 			}
// 			close(readChan)
// 			return content, err

// 		}
// 		return "", errors.New("reqNo not chan")
// 	}

// 	return "", errors.New("reqNo nil obj")
// }

// func (this *TCPConn) SendQuery(content []byte) (reqNo string, size int, err error) {

// 	reqNo = strings.Replace(time.Now().Format("0405.00"), ".", "", 1)
// 	reqNo += fmt.Sprintf("%010d", atomic.AddUint32(&this.reqNo, 1))
// 	if this.reqNo >= 4123456789 {
// 		atomic.StoreUint32(&this.reqNo, 0)
// 	}
// 	query := append([]byte(recv_type_query), []byte(reqNo)...)
// 	size, err = Write(this.Conn, append(query, content...))

// 	ch := make(chan string, 1)
// 	this.chanResult.Store(reqNo, ch)

// 	return reqNo, size, err
// }

// func (this *TCPConn) SendResult(reqNo string, content []byte) (size int, err error) {

// 	query := append([]byte(recv_type_result), []byte(reqNo)...)

// 	return Write(this.Conn, append(query, content...))
// }

// func (this *TCPConn) Close() error {
// 	return this.Conn.Close()
// }

// func Write(conn net.Conn, content []byte) (int, error) {

// 	var lenBuf = make([]byte, 4)
// 	binary.BigEndian.PutUint32(lenBuf, uint32(len(content)))

// 	var buf bytes.Buffer
// 	buf.Write(lenBuf)
// 	buf.Write(content)

// 	return conn.Write(buf.Bytes())
// }

// func Read(conn net.Conn) (string, error) {

// 	var pkg_len_buf = make([]byte, 4)
// 	_, err := conn.Read(pkg_len_buf)
// 	if err != nil {
// 		return "", err
// 	}
// 	var pkg_len = binary.BigEndian.Uint32(pkg_len_buf)
// 	var result = make([]byte, pkg_len)
// 	var read_len int
// 	for {
// 		content := result[read_len:pkg_len]
// 		_len, err := conn.Read(content)
// 		if err != nil {
// 			return "", err
// 		}
// 		read_len += _len
// 		if read_len >= format.Int(pkg_len) {
// 			return string(result), nil
// 		}
// 	}
// }

// type tcpBase struct {
// 	Addr        string
// 	conn        *TCPConn
// 	loopProcess func(conn *TCPConn, async ...bool)
// }

// func (this *tcpBase) Construct() {
// }

// func (this *tcpBase) Conn() *TCPConn {
// 	return this.conn
// }

// func (this *tcpBase) LoopHandler(process func(conn *TCPConn, async ...bool)) {
// 	this.loopProcess = process
// }

// func (this *tcpBase) LoopProcess() {
// 	this.loopProcess(this.conn)
// }

// func (this *tcpBase) Close() error {
// 	return this.conn.Close()
// }

// type TCPClient struct {
// 	tcpBase
// 	Addr string
// }

// func (this *TCPClient) Init() error {

// 	conn, err := net.Dial("tcp", this.Addr)
// 	if err != nil {
// 		return err
// 	}
// 	this.conn = &TCPConn{}
// 	this.conn.Construct()
// 	this.conn.Conn = conn
// 	return nil
// }

// type TCPServ struct {
// 	tcpBase
// 	listener net.Listener
// }

// func (this *TCPServ) Start() error {

// 	if this.Addr == "" {
// 		logger.Run("info", "tcp serv addr empty")
// 		return nil
// 	}
// 	ln, err := net.Listen("tcp", this.Addr)
// 	if err != nil {
// 		return fmt.Errorf("tcp listen error: %v", err)
// 	}
// 	this.listener = ln
// 	logger.Run("info", "tcp serv start", this.Addr)

// 	for {
// 		Conn, err := ln.Accept()
// 		if err != nil {
// 			fmt.Errorf("tcp accept error: %v", err)
// 			break
// 		}
// 		go func(conn net.Conn) {
// 			tcpConn := &TCPConn{Conn: conn}
// 			tcpConn.Construct()
// 			this.loopProcess(tcpConn)
// 		}(Conn)
// 	}
// 	return nil
// }

// func (this *TCPServ) Close() {
// 	this.listener.Close()
// }

// type TCPRequest struct {
// 	content string
// 	Action  string
// 	Params  map[string]interface{}
// 	Cookie  map[string]interface{}
// 	Header  map[string]interface{}
// }

// func (this *TCPRequest) Construct() {

// 	this.Params = make(map[string]interface{})
// 	this.Cookie = make(map[string]interface{})
// 	this.Header = make(map[string]interface{})
// }

// func (this *TCPRequest) Encode() []byte {
// 	data, _ := json.Marshal(this)
// 	return data
// }

// func (this *TCPRequest) Parse(line string) error {

// 	this.content = line
// 	return json.Unmarshal([]byte(line), this)
// }
