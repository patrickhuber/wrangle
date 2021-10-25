package operatingsystem

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
