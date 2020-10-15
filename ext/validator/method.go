package validator

import (
	"fmt"
	"strings"

	"github.com/phper-go/frame/func/object"
)

type Method struct {
	Fields string
	Method func() (errno, errmsg string)
}

func (this *Method) GetFields() string {
	return this.Fields
}

func (this *Method) Check(values map[string]interface{}) (errno, errmsg, field string) {

	if len(this.Fields) != 0 {

		fieldArr := strings.Split(this.Fields, ",")
		for _, attr := range fieldArr {
			if errno, errmsg = this.Method(); errno != "" {
				if strings.Index(errmsg, "%s") >= 0 {
					return errno, fmt.Sprintf(errmsg, attr), attr
				} else {
					return errno, errmsg, attr
				}
			}
			return
		}
	}
	return
}

func (this *Method) CheckObject(obj interface{}) (errno, errmsg, field string) {

	return this.Check(object.Values(obj))
}
