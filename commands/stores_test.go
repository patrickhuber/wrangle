package commands

import (
	"bytes"
	"testing"

	"github.com/patrickhuber/wrangle/config"
	"github.com/patrickhuber/wrangle/ui"
	"github.com/stretchr/testify/require"
)

func TestStoresCommand(t *testing.T) {
	t.Run("CanListStores", func(t *testing.T) {
		r := require.New(t)

		console := ui.NewMemoryConsole()

		cmd := NewStores(console)
		content := `
stores:
- name: one
  type: file
- name: two
  type: credhub
`
		cfg, err := config.DeserializeConfigString(content)
		r.Nil(err)
		err = cmd.Execute(cfg)
		r.Nil(err)

		b, ok := console.Out().(*bytes.Buffer)
		r.True(ok)
		r.NotNil(b)
		r.Equal("name type\n---- ----\none  file\ntwo  credhub\n", b.String())
	})
}
