package archive

type Provider interface {
	Archiver
	Extractor
}
