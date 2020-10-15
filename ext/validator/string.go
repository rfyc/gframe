package validator

import (
	"strings"

	"github.com/phper-go/frame/func/conv"
	"github.com/phper-go/frame/func/object"
)

type String struct {
	Fields     string
	Errno      string
	Errmsg     string
	AllowEmpty bool
	Min        int
	Minmsg     string
	Max        int
	Maxmsg     string
}

func (this *String) GetFields() string {
	return this.Fields
}

func (this *String) Check(values map[string]interface{}) (errno, errmsg, field string) {

	if len(this.Fields) != 0 {
		if this.Errno == "" {
			this.Errno = "119002"
		}
		if this.Errmsg == "" {
			this.Errmsg = "param fail"
		}
		fieldArr := strings.Split(this.Fields, ",")
		for key, val := range values {
			for _, field := range fieldArr {
				if strings.ToLower(field) == strings.ToLower(key) {
					if !this.AllowEmpty && len(conv.String(val)) == 0 {
						return this.Errno, this.Errmsg, field
					}
				}
			}
		}
	}
	return "", "", ""
}

func (this *String) CheckObject(obj interface{}) (errno, errmsg, field string) {

	return this.Check(object.Values(obj))
}
