package store

type FileStore struct {
	Name string
	Path string
}

func (config *FileStore) GetName() string {
	return config.Name
}
