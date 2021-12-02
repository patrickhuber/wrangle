package structio

type Reader interface {
	Read(out interface{}) error
}
