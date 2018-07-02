package packages

import "path/filepath"

// Extract represents an interface for extracting a package of software
type Extract interface {
	Filter() string
	OutFile() string
	OutFolder() string
	OutPath() string
}

type extract struct {
	filter    string
	out       string
	outFolder string
	outPath   string
}

// NewExtract Creates a new extract instance
func NewExtract(filter string, out string, outFolder string) Extract {
	outPath := filepath.Join(out, outFolder)
	return &extract{
		filter:    filter,
		out:       out,
		outFolder: outFolder,
		outPath:   outPath}
}

func (e *extract) Filter() string {
	return e.filter
}

func (e *extract) OutFile() string {
	return e.out
}

func (e *extract) OutFolder() string {
	return e.outFolder
}

func (e *extract) OutPath() string {
	return e.outPath
}
