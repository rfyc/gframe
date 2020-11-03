package output

import (
	"net/http"

	"github.com/phper-go/frame/web/theme"
)

type Output struct {
	Status  string
	Error   string
	Content []byte
	Cookies []*http.Cookie
	Theme   theme.Interface
	Headers map[string]string
}

func (this *Output) Redirect(uri string) {
	this.Headers["Location"] = uri
}
