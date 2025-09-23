package export_test

import (
	"bytes"
	"testing"

	"github.com/patrickhuber/go-cross/console"
	"github.com/patrickhuber/go-shellhook"
	"github.com/patrickhuber/wrangle/internal/envdiff"
	"github.com/patrickhuber/wrangle/internal/export"
	"github.com/stretchr/testify/require"
)

func TestService_Execute_WritesShellCommandsAndDiff(t *testing.T) {
	shells := map[string]shellhook.Shell{
		shellhook.Bash:       shellhook.NewBash(),
		shellhook.Powershell: shellhook.NewPowershell(),
	}

	type testCase struct {
		name       string
		shell      string
		changes    []envdiff.Change
		expectSubs []string
	}

	cases := []testCase{
		{
			name:  "bash_changes",
			shell: shellhook.Bash,
			changes: []envdiff.Change{
				envdiff.Add{Key: "FOO", Value: "bar"},
				envdiff.Update{Key: "FOO", Value: "baz"},
				envdiff.Remove{Key: "BAR"},
			},
			expectSubs: []string{
				"export FOO='bar'",
				"export FOO='baz'",
				"unset BAR",
			},
		},
		{
			name:  "pwsh_changes",
			shell: shellhook.Powershell,
			changes: []envdiff.Change{
				envdiff.Add{Key: "FOO", Value: "bar"},
				envdiff.Update{Key: "FOO", Value: "baz"},
				envdiff.Remove{Key: "BAR"},
			},
			expectSubs: []string{
				"$env:FOO=\"bar\"",
				"$env:FOO=\"baz\"",
				"Remove-Item Env:\\BAR",
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			con := console.NewMemory()
			svc := export.NewService(shells, con)

			err := svc.Execute(tc.shell, tc.changes)
			require.NoError(t, err)

			out := con.Out().(*bytes.Buffer).String()
			require.NotEmpty(t, out)

			for _, sub := range tc.expectSubs {
				require.Contains(t, out, sub)
			}
		})
	}
}

func TestService_Execute_InvalidShell(t *testing.T) {
	shells := map[string]shellhook.Shell{shellhook.Bash: shellhook.NewBash()}
	con := console.NewMemory()
	svc := export.NewService(shells, con)

	err := svc.Execute("zsh", nil)
	require.Error(t, err)
}
