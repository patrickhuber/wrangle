package tasks

// LinkTask represents a symlink
type LinkTask struct {
	Details LinkTaskDetails `yaml:"link" mapstructure:"link"`
}

// LinkTaskDetails contain the details for the link task
type LinkTaskDetails struct {
	Source string `yaml:"source"`
	Alias  string `yaml:"alias"`
}

func (task *LinkTask) Type() string {
	return "link"
}

func (task *LinkTask) Params() map[string]interface{} {
	dictionary := make(map[string]interface{})
	dictionary["source"]= task.Details.Source
	dictionary["alias"]= task.Details.Alias
	return dictionary
}

func NewLinkTask(source string, alias string) Task {
	return &LinkTask{
		Details: LinkTaskDetails{
			Source: source,
			Alias:  alias,
		},
	}
}
