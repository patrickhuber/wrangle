package renderers

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestFactory(t *testing.T) {
	t.Run("WindowsCreatesPowershellRenderer", func(t *testing.T) {
		r := require.New(t)
		platform := "windows"
		shell := ""
		factory := NewFactory()
		renderer, err := factory.Create(shell, platform)
		r.Nil(err)
		r.Equal("powershell", renderer.Shell())
	})

	t.Run("LinuxCreatesBashRenderer", func(t *testing.T) {
		r := require.New(t)
		platform := "linux"
		shell := ""
		factory := NewFactory()
		renderer, err := factory.Create(shell, platform)
		r.Nil(err)
		r.Equal("bash", renderer.Shell())
	})

	t.Run("DarwinCreatesBashRenderer", func(t *testing.T) {
		r := require.New(t)
		platform := "linux"
		shell := ""
		factory := NewFactory()
		renderer, err := factory.Create(shell, platform)
		r.Nil(err)
		r.Equal("bash", renderer.Shell())
	})

	t.Run("BashCreatesBashRenderer", func(t *testing.T) {
		r := require.New(t)
		platform := "linux"
		shell := "bash"
		factory := NewFactory()
		renderer, err := factory.Create(shell, platform)
		r.Nil(err)
		r.Equal("bash", renderer.Shell())
	})

	t.Run("PowershellCreatesPowershellRenderer", func(t *testing.T) {
		r := require.New(t)
		platform := "windows"
		shell := "powershell"
		factory := NewFactory()
		renderer, err := factory.Create(shell, platform)
		r.Nil(err)
		r.Equal("powershell", renderer.Shell())
	})

	t.Run("PlatformIsIgnoredIfShellIsSpecified", func(t *testing.T) {
		r := require.New(t)
		platform := "windows"
		shell := "bash"
		factory := NewFactory()
		renderer, err := factory.Create(shell, platform)
		r.Nil(err)
		r.Equal("bash", renderer.Shell())
	})
}
