package archiver

type Archiver interface {
	Archive(archive string, paths []string) error
	Extract(archive string, destination string, files []string) error
}
