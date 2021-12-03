package operatingsystem

const (
	MockAmd64Architecture = "amd64"
	MockArm64Architecture = "arm64"

	MockWindowsPlatform         = "windows"
	MockWindowsWorkingDirectory = "c:\\working"
	MockWindowsHomeDirectory    = "c:\\users\\fake"
	MockWindowsExecutable       = "c:\\ProgramData\\wrangle\\wrangle.exe"

	MockLinuxPlatform         = "linux"
	MockLinuxWorkingDirectory = "/working"
	MockLinuxHomeDirectory    = "/home/fake"
	MockLinuxExecutable       = "/opt/wrangle/bin/wrangle"

	MockDarwinPlatform         = "darwin"
	MockDarwinHomeDirectory    = MockLinuxHomeDirectory
	MockDarwinWorkingDirectory = MockLinuxWorkingDirectory
	MockDarwinExecutable       = MockLinuxExecutable
)

type mockOS struct {
	workingDirectory string
	platform         string
	architecture     string
	homeDirectory    string
	executable       string
}

type NewMockOS struct {
	WorkingDirectory string
	Platform         string
	Architecture     string
	HomeDirectory    string
	Executable       string
}

// NewMock creates a new OS from the mock OS request
func NewMock(o *NewMockOS) OS {
	return &mockOS{
		executable:       o.Executable,
		workingDirectory: o.WorkingDirectory,
		architecture:     o.Architecture,
		platform:         o.Platform,
		homeDirectory:    o.HomeDirectory,
	}
}

func NewLinuxMock() OS {
	return &mockOS{
		executable:       MockLinuxExecutable,
		workingDirectory: MockLinuxWorkingDirectory,
		platform:         MockLinuxPlatform,
		homeDirectory:    MockLinuxHomeDirectory,
		architecture:     MockAmd64Architecture,
	}
}

func NewDarwinMock() OS {
	return &mockOS{
		executable:       MockDarwinExecutable,
		workingDirectory: MockDarwinWorkingDirectory,
		platform:         MockDarwinPlatform,
		homeDirectory:    MockDarwinHomeDirectory,
		architecture:     MockAmd64Architecture,
	}
}

func NewWindowsMock() OS {
	return &mockOS{
		executable:       MockWindowsExecutable,
		workingDirectory: MockWindowsWorkingDirectory,
		platform:         MockWindowsPlatform,
		homeDirectory:    MockWindowsHomeDirectory,
		architecture:     MockAmd64Architecture,
	}
}

func (o *mockOS) WorkingDirectory() (string, error) {
	return o.workingDirectory, nil
}

func (o *mockOS) Executable() (string, error) {
	return o.executable, nil
}

func (o *mockOS) Platform() string {
	return o.platform
}

func (o *mockOS) Architecture() string {
	return o.architecture
}

func (o *mockOS) Home() string {
	return o.homeDirectory
}
