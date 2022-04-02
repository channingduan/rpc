package validator

type IValidator interface {
	Bind(data interface{}) error
}
