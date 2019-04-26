package packages

type ManifestProvider interface {
	Get(context PackageContext) (*Manifest, error)
	GetInterface(context PackageContext) (interface{}, error)
}
