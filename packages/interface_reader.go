package packages

type InterfaceReader interface {
	Read() (interface{}, error)
}
