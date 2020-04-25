package packages

// ContextProvider defines an interface for a package context provider this
type ContextProvider interface {
	Get(packageName string, packageVerison string) (PackageContext, error)
}
