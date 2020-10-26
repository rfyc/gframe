package route

import (
	"errors"

	"github.com/phper-go/frame/func/conv"
	"github.com/phper-go/frame/interfaces"
	"github.com/phper-go/frame/web/session"
)

func sessionRead(execController interfaces.Controller) error {

	var input = execController.Input()

	if session.Enable > 0 {

		execController.Session().SID = conv.String(input.Cookie[session.Name])
		if len(execController.Session().SID) == 0 {
			return nil
		}

		sessioin_data, err := session.Read(execController.Session().SID)
		if err != nil {
			return errors.New("session read error: " + err.Error())
		}

		for key, val := range sessioin_data {
			execController.Session().Set(key, val)
		}
	}
	return nil
}

func sessionWrite(sess *session.Session) error {
	return session.Write(sess.SID, sess.All())
}
