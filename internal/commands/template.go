package commands

import "github.com/urfave/cli/v2"

// CommandHelpTemplate is the text template for the command help topic.
// cli.go uses text/template to render templates. You can
// render custom help text by setting this variable.
var CommandHelpTemplate = cli.SubcommandHelpTemplate
