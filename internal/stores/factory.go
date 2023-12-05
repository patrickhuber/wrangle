package stores

type Factory interface {
	Name() string
	Create(properties map[string]string) (Store, error)
}
