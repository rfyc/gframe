package validator

import (
	"strings"

	"github.com/phper-go/frame/func/conv"
	"github.com/phper-go/frame/func/object"
)

type Required struct {
	Fields string
	Errno  string
	Errmsg string
}

func (this *Required) GetFields() string {
	return this.Fields
}

func (this *Required) Check(values map[string]interface{}) (errno, errmsg, field string) {

	if len(this.Fields) != 0 {
		if this.Errno == "" {
			this.Errno = errno_required
		}
		if this.Errmsg == "" {
			this.Errmsg = errmsg_required
		}
		fieldArr := strings.Split(this.Fields, ",")
		for key, val := range values {
			for _, field := range fieldArr {
				if strings.ToLower(field) == strings.ToLower(key) {
					if len(conv.String(val)) == 0 {
						return this.Errno, this.Errmsg, field
					}
				}
			}
		}
	}
	return
}

func (this *Required) CheckObject(obj interface{}) (errno, errmsg, field string) {

	return this.Check(object.Values(obj))
}
