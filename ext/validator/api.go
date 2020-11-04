package validator

type Api struct {
	field  string
	errno  string
	errmsg string
}

func (this *Api) GetErrors() (errno, errmsg, field string) {

	return this.errno, this.errmsg, this.field
}

func (this *Api) SetErrors(errno, errmsg, field string) {
	this.errno = errno
	this.errmsg = errmsg
	this.field = field
}

func (this *Api) Rules() Rules {

	return Rules{}
}
