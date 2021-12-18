package archive

type Archiver interface {
	Archive(archive string, paths ...string) error
}
