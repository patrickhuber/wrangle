package archive

type Extractor interface {
	Extract(archive string, destination string, files ...string) error
}
