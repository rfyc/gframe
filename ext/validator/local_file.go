package validator

import (
	"strings"

	"github.com/phper-go/frame/func/conv"
	"github.com/phper-go/frame/func/file"
	"github.com/phper-go/frame/func/object"
)

type LocalFile struct {
	Fields     string
	AllowEmpty bool
	Errno      string
	Errmsg     string
}

func (this *LocalFile) GetFields() string {
	return this.Fields
}

func (this *LocalFile) Check(values map[string]interface{}) (errno, errmsg, field string) {

	if len(this.Fields) != 0 {
		if this.Errno == "" {
			this.Errno = errno_local_file
		}
		if this.Errmsg == "" {
			this.Errmsg = errmsg_local_file
		}
		fieldArr := strings.Split(this.Fields, ",")
		for key, val := range values {
			for _, f := range fieldArr {
				if strings.ToLower(f) == strings.ToLower(key) {
					if !this.AllowEmpty && !file.IsFile(conv.String(val)) {
						return this.Errno, this.Errmsg, f
					}
				}
			}
		}
	}
	return
}

func (this *LocalFile) CheckObject(obj interface{}) (errno, errmsg, field string) {

	return this.Check(object.Values(obj))
}
