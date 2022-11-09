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
	Steps        []ManifestStep
}

type ManifestStep struct {
	Action string
	With   map[string]any
}
