package packages

type Manifest struct {
	Package *ManifestPackage `yaml:"package" json:"package"`
}

type ManifestPackage struct {
	Name    string
	Version string
	Targets []*ManifestTarget
}

type ManifestTarget struct {
	Platform     string
	Architecture string
	Tasks        []map[string]map[string]string
}
