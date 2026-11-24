package commands

type StoreAddOptions struct {
	Name         string   `position:"0"`
	Type         string   `flag:"type"`
	Properties   []string `flag:"property"`
	Include      []string `flag:"include"`
	Exclude      []string `flag:"exclude"`
	Dependencies []string `flag:"dependency"`
	Scope        string   `flag:"scope"`
}
