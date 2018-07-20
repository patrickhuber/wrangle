package renderers

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMain(m *testing.M) {
	os.Unsetenv("PSModulePath")
	m.Run()
}

func TestFactory(t *testing.T) {
	t.Run("WindowsCreatesPowershellRenderer", func(t *testing.T) {
		r := require.New(t)
		platform := "windows"
		shell := ""
		factory := NewFactory(platform)
		renderer, err := factory.Create(shell)
		r.Nil(err)
		r.Equal("powershell", renderer.Shell())
	})

	t.Run("LinuxCreatesBashRenderer", func(t *testing.T) {
		r := require.New(t)
		platform := "linux"
		shell := ""
		factory := NewFactory(platform)
		renderer, err := factory.Create(shell)
		r.Nil(err)
		r.Equal("bash", renderer.Shell())
	})

	t.Run("DarwinCreatesBashRenderer", func(t *testing.T) {
		r := require.New(t)
		platform := "linux"
		shell := ""
		factory := NewFactory(platform)
		renderer, err := factory.Create(shell)
		r.Nil(err)
		r.Equal("bash", renderer.Shell())
	})

	t.Run("LinuxWithPsModuelPathCreatesPowershellRenderer", func(t *testing.T) {
		r := require.New(t)
		platform := "linux"
		shell := ""
		err := os.Setenv("PSModulePath", "test")
		r.Nil(err)

		factory := NewFactory(platform)
		renderer, err := factory.Create(shell)
		r.Nil(err)

		err = os.Unsetenv("PSModulePath")
		r.Nil(err)

		r.Equal("powershell", renderer.Shell())
	})

	t.Run("DarwinWithPsModuelPathCreatesPowershellRenderer", func(t *testing.T) {
		r := require.New(t)
		platform := "darwin"
		shell := ""
		err := os.Setenv("PSModulePath", "test")
		r.Nil(err)

		factory := NewFactory(platform)
		renderer, err := factory.Create(shell)
		r.Nil(err)

		err = os.Unsetenv("PSModulePath")
		r.Nil(err)

		r.Equal("powershell", renderer.Shell())
	})

	t.Run("BashCreatesBashRenderer", func(t *testing.T) {
		r := require.New(t)
		platform := "linux"
		shell := "bash"
		factory := NewFactory(platform)
		renderer, err := factory.Create(shell)
		r.Nil(err)
		r.Equal("bash", renderer.Shell())
	})

	t.Run("PowershellCreatesPowershellRenderer", func(t *testing.T) {
		r := require.New(t)
		platform := "windows"
		shell := "powershell"
		factory := NewFactory(platform)
		renderer, err := factory.Create(shell)
		r.Nil(err)
		r.Equal("powershell", renderer.Shell())
	})

	t.Run("PlatformIsIgnoredIfShellIsSpecified", func(t *testing.T) {
		r := require.New(t)
		platform := "windows"
		shell := "bash"
		factory := NewFactory(platform)
		renderer, err := factory.Create(shell)
		r.Nil(err)
		r.Equal("bash", renderer.Shell())
	})
}
