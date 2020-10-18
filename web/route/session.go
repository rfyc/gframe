package route

import (
	"errors"
	"net/http"
	"time"

	"github.com/phper-go/frame/func/conv"
	"github.com/phper-go/frame/web/input"
	"github.com/phper-go/frame/web/output"
	"github.com/phper-go/frame/web/session"
)

func sessionRead(input *input.Input, output *output.Output) error {

	if session.Enable > 0 {
		session_id := conv.String(input.Cookie[session.Name])
		if len(session_id) == 0 {
			session_id = session.ID()
			output.Cookies = append(output.Cookies, &http.Cookie{
				Name:    session.Name,
				Value:   session_id,
				Expires: time.Unix(time.Now().Unix()+int64(session.LifeTime), 0),
			})
			input.Cookie[session.Name] = session_id
			return nil
		}
		sessioin_data, err := session.Read(session_id)
		if err != nil {
			return errors.New("session read error: " + err.Error())
		}
		input.Session = sessioin_data
	}
	return nil
}

func sessionWrite(input *input.Input) error {
	session_id := conv.String(input.Cookie[session.Name])
	return session.Write(session_id, input.Session)
}
