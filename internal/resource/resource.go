package resource

// Resource is a generic way to fetch versions from a given endpoint
type Resource interface {
	Type() string
	Name() string
	Versions() ([]string, error)
}
