package global

const (
	// PackagePathKey defines the environment variable for the package cache
	PackagePathKey = "WRANGLE_PACKAGES"

	// BinPathKey defines the environment variable for the bin directory where package links are installed
	BinPathKey = "WRANGLE_BIN"

	// ConfigFileKey defines the environment variable for the current configuration
	ConfigFileKey = "WRANGLE_CONFIG"

	// RootPathKey defines the root directory wrangle uses. If specified, defaults for WRANGLE_PACKAGES, WRANGLE_BIN can be infered
	RootPathKey = "WRANGLE_ROOT"
)
