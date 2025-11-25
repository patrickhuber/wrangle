package stores

type Factory interface {
	Name() string
	Create(properties map[string]any) (Store, error)
}
