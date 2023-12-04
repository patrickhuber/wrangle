package structio

// Writer defines a structured writer
type Writer interface {
	Write(any) error
}
