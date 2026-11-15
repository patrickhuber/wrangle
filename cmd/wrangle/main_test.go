package main

import (
	"bytes"
	"strings"
	"testing"

	"github.com/urfave/cli/v2"
)

func TestInstallCommandIsListedAtRoot(t *testing.T) {
	var output bytes.Buffer
	app := &cli.App{
		Name:     "wrangle",
		Commands: rootCommands,
		Writer:   &output,
	}

	if err := app.Run([]string{"wrangle", "--help"}); err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(output.String(), "install") {
		t.Fatalf("root help does not list install command:\n%s", output.String())
	}

	err := app.Run([]string{"wrangle", "install"})
	if err == nil || err.Error() != "package name is required" {
		t.Fatalf("install command did not validate its package argument: %v", err)
	}
}
