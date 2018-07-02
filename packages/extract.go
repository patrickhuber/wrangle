package packages

// Extract represents an interface for extracting a package of software
type Extract interface {
	Filter() string
	Out() string
	OutFolder() string
}

type extract struct {
	filter    string
	out       string
	outFolder string
}

// NewExtract Creates a new extract instance
func NewExtract(filter string, out string, outFolder string) Extract {
	return &extract{filter: filter, out: out, outFolder: outFolder}
}

func (e *extract) Filter() string {
	return e.filter
}

func (e *extract) Out() string {
	return e.out
}

func (e *extract) OutFolder() string {
	return e.outFolder
}
