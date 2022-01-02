package packages

type Manifest struct {
	Package *ManifestPackage `yaml:"package", json"package"`
}

type ManifestPackage struct {
	Name    string
	Version string
	Targets []*Target
}
