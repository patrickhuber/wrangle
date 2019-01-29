package commands

import (
	"fmt"

	"github.com/urfave/cli"
)

func setCommandCustomHelpTemplateWithGlobalOptions(app *cli.App, command *cli.Command) {
	template := cli.CommandHelpTemplate
	template += fmt.Sprintln()
	template += fmt.Sprintf("GLOBAL OPTIONS:")
	template += fmt.Sprintln()
	for _, globalFlag := range app.Flags {
		template += fmt.Sprintf("\t%v", globalFlag)
		template += fmt.Sprintln()
	}
	command.CustomHelpTemplate = template
	command.UsageText = fmt.Sprintf("%s [ global options ] %s", app.HelpName, command.Name)
}
