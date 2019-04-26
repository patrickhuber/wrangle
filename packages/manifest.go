package packages

// Manifest represents a versioned artifiact
type Manifest struct {
	Name    string   `yaml:"name"`
	Version string   `yaml:"version"`
	Targets []Target `yaml:"targets"`
}

// Target repesents an install target
type Target struct {
	Platform     string        `yaml:"platform"`
	Architecture string        `yaml:"architecture"`
	Tasks        []interface{} `yaml:"tasks"`
}
