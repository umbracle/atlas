package cmd

import (
	"github.com/mitchellh/cli"
)

// ProvidersCommand is the command to show the version of the agent
type ProvidersCommand struct {
	UI cli.Ui
}

// Help implements the cli.Command interface
func (c *ProvidersCommand) Help() string {
	return `Usage: ensemble version
	
  Display the Ensemble version`
}

// Synopsis implements the cli.Command interface
func (c *ProvidersCommand) Synopsis() string {
	return "Display the Ensemble version"
}

// Run implements the cli.Command interface
func (c *ProvidersCommand) Run(args []string) int {
	return 0
}
