package stores

type Store interface {
	Readable
	Writable
}

// Readable store is the read interface for Store
type Readable interface {
	Get(k Key) (any, bool, error)
	List() ([]Key, error)
}

// Writeable store is the write interface for Store
type Writable interface {
	Set(k Key, v any) error
}
