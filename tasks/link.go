package tasks

// LinkTask represents a symlink
type LinkTask struct {
	Details LinkTaskDetails `yaml:"link"`
}

// LinkTaskDetails contain the details for the link task
type LinkTaskDetails struct {
	Source      string `yaml:"source"`
	Destination string `yaml:"destination"`
}
