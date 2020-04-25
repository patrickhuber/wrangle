package feed

// Writer writes the packages
type Writer interface {
	Write(packages []*Package) error
}
