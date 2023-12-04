package stores

type Store interface {
	Get(k Key) (any, bool, error)
	List() ([]Key, error)
}
