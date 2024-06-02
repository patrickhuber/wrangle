package stores

type Store interface {
	Get(k Key) (any, bool, error)
	Set(k Key, v any) error
	List() ([]Key, error)
}
