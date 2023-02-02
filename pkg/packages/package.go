package packages

// Package defines a software package and how it is installed
type Package struct {
	Name     string
	Versions []*Version
}

// Version returns the version of a package
type Version struct {
	Version  string
	Manifest *Manifest
}
