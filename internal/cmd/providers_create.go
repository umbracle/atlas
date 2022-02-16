package cmd

// ProvidersCreateCommand is the command to show the version of the agent
type ProvidersCreateCommand struct {
	*Meta
}

// Help implements the cli.Command interface
func (c *ProvidersCreateCommand) Help() string {
	return `Usage: ensemble version
	
  Display the Ensemble version`
}

// Synopsis implements the cli.Command interface
func (c *ProvidersCreateCommand) Synopsis() string {
	return "Display the Ensemble version"
}

// Run implements the cli.Command interface
func (c *ProvidersCreateCommand) Run(args []string) int {
	return 0
}
