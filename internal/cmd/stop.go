package cmd

import (
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/umbracle/atlas/internal/state"
)

// StopCommand is the command to show the version of the agent
type StopCommand struct {
	UI cli.Ui
}

// Help implements the cli.Command interface
func (c *StopCommand) Help() string {
	return `Usage: ensemble version
	
  Display the Ensemble version`
}

// Synopsis implements the cli.Command interface
func (c *StopCommand) Synopsis() string {
	return "Display the Ensemble version"
}

// Run implements the cli.Command interface
func (c *StopCommand) Run(args []string) int {

	arg := args[0]
	fmt.Println(arg)

	state, err := state.NewState("./state.db")
	if err != nil {
		panic(err)
	}

	nodes, err := state.ListNodes()
	if err != nil {
		panic(err)
	}
	fmt.Println(nodes)

	return 0
}
