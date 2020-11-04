package interfaces

import (
	"net"

	"github.com/phper-go/frame/web/input"
	"github.com/phper-go/frame/web/output"
	"github.com/phper-go/frame/web/session"
)

type Config interface {
	EnvName() string
	EnvFile() string
	DefaultFile() string
	Load(app App, file ...string) error
}

type App interface {
	Construct()
	Init() error
	Start() error
	Wrap(func())
	Wait()
}

type DBModel interface {
	PrimaryKey() string
	Table() string
}

type Command interface {
	Construct(runApp App)
	Init() error
	LoadApp() error
	InitApp() error
	Prepare() error
	Run()
	End()
}

type Controller interface {
	Construct(controllerName, actionName string)
	Prepare() bool
	End()
	Input() *input.Input
	Output() *output.Output
	Session() *session.Session
}

type Action interface {
	Controller
	Run()
}

type TCPHandler interface {
	ServeTCP(net.Conn)
}
