package chttp

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/phper-go/frame/func/conv"
)

type CHttp struct {
	Method   string
	URI      string
	Query    url.Values
	Headers  map[string]string
	Cookies  []*http.Cookie
	Client   *http.Client
	Request  *http.Request
	Response *http.Response
}

func (this *CHttp) Do(binds ...interface{}) (content []byte, err error) {

	if content, err = this.do(); err != nil {
		return content, err
	} else if 200 != this.Response.StatusCode {
		return content, errors.New("http code error: " + this.Response.Status)
	} else if len(content) == 0 {
		return content, errors.New("content empty error")
	} else {
		if len(binds) > 0 {
			for _, bind := range binds {
				if err = json.Unmarshal(content, bind); err != nil {
					return content, errors.New("bind error: " + err.Error())
				}
			}
		}
		return content, err
	}
}

func (this *CHttp) do() (content []byte, err error) {

	//request new

	reader := strings.NewReader(this.Query.Encode())
	if request, err := http.NewRequest(this.Method, this.URI, reader); err != nil {
		return []byte{}, errors.New("request new error: " + err.Error())
	} else {
		//request set
		if this.Headers != nil && len(this.Headers) > 0 {
			for key, value := range this.Headers {
				request.Header.Add(key, value)
			}
		}
		if strings.ToUpper(this.Method) == "POST" {
			request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		} else {
			request.URL.RawQuery = this.Query.Encode()
		}
		for _, cookie := range this.Cookies {
			request.AddCookie(cookie)
		}
		//client do
		if this.Response, err = this.Client.Do(request); err != nil {
			return []byte{}, errors.New("request do error: " + err.Error())
		}
		if body, err := ioutil.ReadAll(this.Response.Body); err != nil {
			return []byte{}, errors.New("body read error: " + err.Error())
		} else if err = this.Response.Body.Close(); err != nil {
			return []byte{}, errors.New("body close error: " + err.Error())
		} else {
			return body, nil
		}
	}

}

func (this *CHttp) SetHeader(key, value string) *CHttp {
	this.Headers[key] = value
	return this
}

func (this *CHttp) SetCookie(cookie *http.Cookie) *CHttp {
	this.Cookies = append(this.Cookies, cookie)
	return this
}

func New(timeout int, method, uri string, query map[string]interface{}) *CHttp {
	client := &CHttp{
		Method:   strings.ToUpper(method),
		URI:      uri,
		Query:    url.Values{},
		Headers:  make(map[string]string),
		Response: &http.Response{},
		Client: &http.Client{
			Timeout: time.Duration(timeout) * time.Millisecond,
		},
	}
	for key, value := range query {
		client.Query.Set(key, conv.String(value))
	}
	return client
}
