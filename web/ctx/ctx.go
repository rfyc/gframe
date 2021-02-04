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
	Server     *Server
	Theme      theme.Interface
}

func NewCtx(parent context.Context) *Ctx {

	var ctx = &Ctx{
		Context: parent,
		Input:   NewInput(),
		Output:  NewOutput(),
		Session: &session.Session{},
		Server:  &Server{},
		Theme:   &theme.Theme{},
	}

	ctx.Theme.Construct()
	ctx.Theme.SetBuffer(&ctx.Output.Content)

	return ctx
}

func NewCtxHTTP(request *http.Request, response http.ResponseWriter) *Ctx {

	var ctx = &Ctx{
		Context: request.Context(),
		Input:   NewInputHTTP(request),
		Output:  NewOutput(),
		Session: &session.Session{},
		Server:  NewServerHTTP(request),
		Theme:   &theme.Theme{},
	}

	ctx.Theme.Construct()
	ctx.Theme.SetBuffer(&ctx.Output.Content)

	return ctx
}

func NewCtxTCP(parent context.Context) *Ctx {
	return NewCtx(parent)
}
