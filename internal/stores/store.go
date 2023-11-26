package stores

type Store interface {
	Get(k Key) (any, error)
	Lookup(k Key) (any, bool, error)
}
