package tasks

// TaskContext defines the context for the given package
type TaskContext interface {
	Root() string
	Bin() string
	PackagesRoot() string
	PackagePath() string
	PackageVersionPath() string
	PackageVersionManifestPath() string
}
