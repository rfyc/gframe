package validator

import (
	"github.com/phper-go/frame/func/object"
)

var (
	errno_default     = "1300"
	errno_required    = "1301"
	errno_local_file  = "1302"
	errmsg_default    = "param fail"
	errmsg_required   = "required"
	errmsg_local_file = "not exist"
)

type ApiInterface interface {
	GetErrors() (errno, errmsg, field string)
	SetErrors(errno, errmsg, field string)
}

type Rules []Rule

type Interface interface {
	Rules() Rules
}

type Rule interface {
	GetFields() string
	CheckObject(obj interface{}) (errno, errmsg, field string)
	Check(values map[string]interface{}) (errno, errmsg, field string)
}

func Merge(rules1, rules2 Rules) Rules {
	for _, rule := range rules2 {
		rules1 = append(rules1, rule)
	}
	return rules1
}

func Check(validator Interface) (errno, errmsg, field string) {

	values := object.Values(validator)

	rules := validator.Rules()

	for _, rule := range rules {

		errno, errmsg, field = rule.Check(values)
		if errno != "" || field != "" {
			if errno == "" {
				errno = errno_default
			}
			if errmsg == "" {
				errmsg = errmsg_default
			}

			if api, ok := validator.(ApiInterface); ok {
				api.SetErrors(errno, errmsg, field)
			}

			return
		}
	}
	return
}
