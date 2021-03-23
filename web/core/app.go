package core

import (
	"errors"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
	"syscall"

	"github.com/phper-go/frame/web/route"

	"github.com/phper-go/frame/logger"

	"github.com/phper-go/frame/connect"
	"github.com/phper-go/frame/ext/signal"
	"github.com/phper-go/frame/func/conv"
)

// type setting struct {
// 	// Session session.Config
// }

type App struct {
	Debug         int
	DefaultAction string
	Logs          string
	HTTPServer    *connect.HTTPServer
	TCPServer     *connect.TCPServer
	waitGroup     sync.WaitGroup
}

func (this *App) Construct() {

	route.Debug = &this.Debug
	route.DefaultAction = &this.DefaultAction

	this.HTTPServer = &connect.HTTPServer{}
	this.HTTPServer.Construct()
	this.HTTPServer.Handle("/", &route.Handler{})

	// this.TCPServer = &connect.TCPServer{}
	// this.TCPServer.Construct()
	// this.TCPServer.Handle("/", &route.TCP{})
}

func (this *App) initLogger() error {

	if this.Logs == "" {
		if execPath, err := os.Executable(); err != nil {
			return errors.New("execpath - " + err.Error())
		} else {
			this.Logs = filepath.Dir(execPath) + "/logs"
		}
		if err := os.MkdirAll(this.Logs, os.ModePerm); err != nil {
			return errors.New("mkdir - " + err.Error())
		}
		if err := syscall.Access(this.Logs, syscall.O_RDWR); err != nil {
			return errors.New("writable - " + err.Error())
		}
	}
	logger.Path = this.Logs
	return nil
}

func (this *App) Init() error {

	if err := this.initLogger(); err != nil {
		return errors.New("logger - " + err.Error())
	}

	if err := this.HTTPServer.Init(); err != nil {
		return errors.New("httpserv - " + err.Error())
	}

	if err := this.TCPServer.Init(); err != nil {
		return errors.New("tcpserv - " + err.Error())
	}

	return nil
}

func (this *App) Start() error {

	this.Wrap(func() {
		if this.HTTPServer.Address != "" {
			logger.Format("listen", this.HTTPServer.Address, "http").Echo(os.Stderr)
			err := this.HTTPServer.Start()
			logger.Format("listen end http -", err.Error()).Echo(os.Stderr)
		} else {
			logger.Format("addr empty http").Echo(os.Stderr)
		}
	})

	// this.Wrap(func() {
	// 	if this.TCPServer.Address != "" {
	// 		logger.Format("listen", this.TCPServer.Address, "tcp").Echo(os.Stderr)
	// 		err := this.TCPServer.Start()
	// 		logger.Format("listen end tcp -", err.Error()).Echo(os.Stderr)
	// 	} else {
	// 		logger.Format("addr empty tcp").Echo(os.Stderr)
	// 	}
	// })

	// this.Wrap(func() {
	// 	ppid := os.Getppid()
	// 	if os.Getenv("graceful") == "on" && ppid > 1 {
	// 		proc, err := os.FindProcess(ppid)
	// 		if err == nil {
	// 			logger.Run("info", "graceful stop", ppid)
	// 			proc.Signal(syscall.SIGTERM)
	// 		}
	// 	}
	// })

	this.signalHandler()

	return nil
}

func (this *App) Wrap(method func()) {
	this.waitGroup.Add(1)
	go func() {
		method()
		this.waitGroup.Done()
	}()
}

func (this *App) Wait() {
	this.waitGroup.Wait()
}

func (this *App) Destruct() {

}

func (this *App) extraFiles() []*os.File {

	var files []*os.File
	var num = 3
	httpListener := this.HTTPServer.Listener()
	if httpListener != nil {
		if listener, ok := httpListener.(*net.TCPListener); ok {
			if fd, err := listener.File(); err == nil {
				files = append(files, fd)
				os.Setenv(this.HTTPServer.Address, conv.String(num))
				num++
			}
		}
	}
	return files
}

func (this *App) signalHandler() {

	signal.Register(syscall.SIGUSR2, func(sig os.Signal) {

		os.Setenv("graceful", "on")
		Cmd := exec.Command(os.Args[0], os.Args[1:]...)
		Cmd.Env = os.Environ()
		Cmd.ExtraFiles = this.extraFiles()
		Cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		Cmd.Start()

	})

	signal.Register(syscall.SIGTERM, func(sig os.Signal) {

		this.HTTPServer.Stop()

	})

}
