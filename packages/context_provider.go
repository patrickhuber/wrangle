package packages

type ContextProvider interface {
	Get(packageName string, packageVerison string) (PackageContext, error)
}
