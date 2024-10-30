package packages

import (
	"github.com/patrickhuber/go-xplat/arch"
	"github.com/patrickhuber/go-xplat/platform"
)

type Manifest struct {
	Package *ManifestPackage `yaml:"package" json:"package"`
}

type ManifestPackage struct {
	Name    string
	Version string
	Targets []*ManifestTarget
}

type ManifestTarget struct {
	Platform     platform.Platform
	Architecture arch.Arch
	Executables  []string
	Steps        []*ManifestStep
}

type ManifestStep struct {
	Action string
	With   map[string]any
}
