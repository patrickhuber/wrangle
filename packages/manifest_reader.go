package packages

type ManifestReader interface {
	Read() (*Manifest, error)
}
