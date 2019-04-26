package packages

type ManifestWriter interface {
	Write(manifest *Manifest) error
}
