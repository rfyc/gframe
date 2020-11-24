package ctx

import (
	"context"

	"github.com/phper-go/frame/web/theme"
)

type Context struct {
	Action     string
	Context    context.Context
	Controller string
	Input      *Input
	Output     *Output
	Theme      theme.Interface
}

func NewCtx(parent context.Context) *Context {

	return &Context{
		ctx:    parent,
		Input:  &Input{},
		Output: &Output{},
	}
}

func NewCtxHTTP(parent context.Context, response http.ResponseWriter, request *http.Request) *Context {

	var ctx = &Context{
		ctx:    parent,
		Input:  NewInputHTTP(request),
		Output: &Output{},
	}

	return ctx
}

func NewCtxTCP(parent context.Context) *Context {
	return New(parent)
}
