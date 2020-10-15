package core

import (
	"github.com/phper-go/frame/interfaces"
)

type Api struct {
	errno  string
	errmsg string
	field  string
}

func (this *Api) Rules() interfaces.ValidatorRules {
	return interfaces.ValidatorRules{}
}

func (this *Api) SetErrors(errno, errmsg, field string) {
	this.errno = errno
	this.errmsg = errmsg
	this.field = field
}

func (this *Api) GetErrors() (errno, errmsg, field string) {
	return this.errno, this.errmsg, this.field
}
