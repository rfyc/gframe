package interfaces

import (
	"net"

	"github.com/phper-go/frame/web/input"
	"github.com/phper-go/frame/web/output"
)

type Api interface {
	Validator
	SetErrors(errno, errmsg, field string)
	GetErrors() (errno, errmsg, field string)
}

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
	Rules() ValidatorRules
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
}

type Action interface {
	Controller
	Run()
}

type TCPHandler interface {
	ServeTCP(net.Conn)
}

type Validator interface {
	Rules() ValidatorRules
}

type ValidatorRule interface {
	GetFields() string
	Check(values map[string]interface{}) (errno, errmsg, field string)
	CheckObject(obj interface{}) (errno, errmsg, field string)
}

type ValidatorRules []ValidatorRule
