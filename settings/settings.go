package settings

type Settings struct {
	Feeds []string `yaml:"feeds"`
	Paths *Paths   `yaml:"paths"`
}

type Paths struct {
	Root     string `yaml:"root"`
	Bin      string `yaml:"bin"`
	Packages string `yaml:"packages"`
}
