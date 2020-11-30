package core

import (
	"github.com/phper-go/frame/web/ctx"
)

type Controller struct {
	ctx *ctx.Ctx
}

func (this *Controller) Ctx() *ctx.Ctx {
	if this.ctx == nil {
		this.ctx = &ctx.Ctx{}
	}
	return this.ctx
}

func (this *Controller) Prepare() bool {
	return true
}

func (this *Controller) End() {

}
