package settings

const (
	// DefaultWindowsRoot is the default root for wrangle on windows
	DefaultWindowsRoot = "c:/tools/wrangle"

	// DefaultWindowsBin is the default windows binary location on windows
	DefaultWindowsBin = DefaultWindowsRoot + "/bin"

	// DefaultWindowsPackages is the default package root on windows
	DefaultWindowsPackages = DefaultWindowsRoot + "/packages"

	// DefaultNixRoot is the default root for *nix systems
	DefaultNixRoot = "/opt/wrangle"

	// DefaultNixBin is the default binary location on *nix systems
	DefaultNixBin = DefaultNixRoot + "/bin"

	// DefaultNixPackages is the default package location on *nix systems
	DefaultNixPackages = DefaultNixRoot + "/packages"
)
