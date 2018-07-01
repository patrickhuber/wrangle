package packages

// Extract represents an interface for extracting a package of software
type Extract interface {
	Filter() string
	Out() string
}

type extract struct {
	filter string
	out    string
}

// NewExtract Creates a new extract instance
func NewExtract(filter string, out string) Extract {
	return &extract{filter: filter, out: out}
}

func (e *extract) Filter() string {
	return e.filter
}

func (e *extract) Out() string {
	return e.out
}
