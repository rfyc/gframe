package validator

import (
	"github.com/phper-go/frame/func/object"
	"github.com/phper-go/frame/interfaces"
)

var (
	errno_default     = "1300"
	errno_required    = "1301"
	errno_local_file  = "1302"
	errmsg_default    = "param fail"
	errmsg_required   = "required"
	errmsg_local_file = "not exist"
)

func MergeRules(rules, rules1 interfaces.ValidatorRules) interfaces.ValidatorRules {
	for _, rule := range rules1 {
		rules = append(rules, rule)
	}
	return rules
}

func Check(obj interfaces.Validator) (errno, errmsg, field string) {

	values := object.Values(obj)

	rules := obj.Rules()

	for _, rule := range rules {

		errno, errmsg, field = rule.Check(values)
		if errno != "" || field != "" {
			if errno == "" {
				errno = "1300"
			}
			if errmsg == "" {
				errmsg = "param fail"
			}
			return
		}
	}
	return
}

func CheckApi(obj interfaces.Api) bool {
	errno, errmsg, field := Check(obj)
	if field != "" || errno != "" {
		obj.SetErrors(errno, errmsg, field)
		return false
	}
	return true
}
