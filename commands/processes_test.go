package commands

import (
	"bytes"
	"testing"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/stretchr/testify/require"
)

func TestProcessesCommand(t *testing.T) {
	t.Run("CanListProcesses", func(t *testing.T) {
		r := require.New(t)

		console := ui.NewMemoryConsole()

		cmd := NewListProcesses(console)
		content := `
environments:
- name: lab  
  processes:
  - name: go 
  - name: echo
- name: prod  
  processes:
  - name: wrangle
  - name: dangle
`
		cfg, err := config.SerializeString(content)
		r.Nil(err)
		err = cmd.Execute(cfg, "lab")
		r.Nil(err)

		b, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(b)
		r.Equal("name\n----\ngo\necho\n", b.String())
	})
}