package cmd

import (
	"fmt"

	"github.com/mitchellh/cli"
	"github.com/umbracle/atlas/internal/state"
)

// ListCommand is the command to show the version of the agent
type ListCommand struct {
	UI cli.Ui
}

// Help implements the cli.Command interface
func (c *ListCommand) Help() string {
	return `Usage: ensemble version
	
  Display the Ensemble version`
}

// Synopsis implements the cli.Command interface
func (c *ListCommand) Synopsis() string {
	return "Display the Ensemble version"
}

// Run implements the cli.Command interface
func (c *ListCommand) Run(args []string) int {

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
