package ctx

import (
	"net/http"
)

type Output struct {
	Status  string
	Error   string
	Content []byte
	Cookies []*http.Cookie
	Headers map[string]string
}

func (this *Output) Redirect(uri string) {
	this.Headers["Location"] = uri
}

func NewOutput() *Output {

	var output = &Output{
		Headers: make(map[string]string),
		Status:  "200",
		Content: []byte{},
	}

	return output

}
