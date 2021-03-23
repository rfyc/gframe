package interfaces

import (
	"net"

	"github.com/phper-go/frame/web/ctx"
)

type App interface {
	Construct()
	Init() error
	Start() error
	Wrap(func())
	Wait()
}

type Config interface {
	EnvName() string
	EnvFile() string
	DefaultFile() string
	Load(app App, file ...string) error
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
	Ctx() *ctx.Ctx
	Prepare() bool
	End()
}

type Action interface {
	Controller
	Run()
}

type TCPHandler interface {
	ServeTCP(net.Conn)
}
