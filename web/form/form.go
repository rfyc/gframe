package form

type FormInterface interface {
	HasError() bool
	GetError() string
	Elements()
	Buttons()
}

type Form struct {
	builder Builder
}

func (this *Form) HasError() bool {

	return true
}

func (this *Form) GetError() string {

	return ""
}
