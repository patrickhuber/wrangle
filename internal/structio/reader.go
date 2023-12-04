package structio

type Reader interface {
	Read(out any) error
}
