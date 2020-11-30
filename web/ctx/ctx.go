package ctx

import (
	"context"
	"net/http"

	"github.com/phper-go/frame/web/session"
	"github.com/phper-go/frame/web/theme"
)

type Ctx struct {
	Action     string
	Context    context.Context
	Controller string
	Input      *Input
	Output     *Output
	Session    *session.Session
	Theme      theme.Interface
}

func NewCtx(parent context.Context) *Ctx {

	var ctx = &Ctx{
		Context: parent,
		Input:   NewInput(),
		Output:  NewOutput(),
		Session: &session.Session{},
		Theme:   &theme.Theme{},
	}

	ctx.Theme.Construct()
	ctx.Theme.SetBuffer(&ctx.Output.Content)

	return ctx
}

func NewCtxHTTP(response http.ResponseWriter, request *http.Request) *Ctx {

	var ctx = &Ctx{
		Context: request.Context(),
		Input:   NewInputHTTP(request),
		Output:  NewOutput(),
		Session: &session.Session{},
		Theme:   &theme.Theme{},
	}

	ctx.Theme.Construct()
	ctx.Theme.SetBuffer(&ctx.Output.Content)

	return ctx
}

func NewCtxTCP(parent context.Context) *Ctx {
	return NewCtx(parent)
}
