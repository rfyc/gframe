package output

import (
	"net/http"
)

type Output struct {
	Status  int
	Error   string
	Content []byte
	Cookies []*http.Cookie
	Headers map[string]string
}

func (this *Output) Redirect(uri string) {
	this.Headers["Location"] = uri
	this.Status = http.StatusFound
}
