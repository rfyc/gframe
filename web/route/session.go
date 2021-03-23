package route

import (
	"errors"

	"github.com/phper-go/frame/func/conv"
	"github.com/phper-go/frame/web/ctx"
	"github.com/phper-go/frame/web/session"
)

func sessionRead(Ctx *ctx.Ctx) error {

	var input = Ctx.Input
	var sess = Ctx.Session
	if session.Enable > 0 {

		sess.SID = conv.String(input.Cookie[session.Name])
		if len(sess.SID) == 0 {
			return nil
		}

		sessioin_data, err := session.Read(sess.SID)
		if err != nil {
			return errors.New("session read error: " + err.Error())
		}

		for key, val := range sessioin_data {
			sess.Set(key, val)
		}
	}
	return nil
}

func sessionWrite(sess *session.Session) error {
	return session.Write(sess.SID, sess.All())
}
