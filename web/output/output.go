package output

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
