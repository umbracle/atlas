package cmd

import (
	"github.com/mitchellh/cli"
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
	panic("TODO")

	/*
		arg := args[0]
		fmt.Println(arg)

		state, err := state.NewState("./state.db")
		if err != nil {
			panic(err)
		}

		d, err := runtime.NewDocker()
		if err != nil {
			panic(err)
		}

		nodes, err := state.ListNodes()
		if err != nil {
			panic(err)
		}
		// find the node
		var node *proto.Node
		for _, nodeTarget := range nodes {
			if nodeTarget.Id == arg {
				node = nodeTarget
			}
		}
		if node == nil {
			c.UI.Error("node not found")
			return 1
		}

		fmt.Println(node)
		d.Stop(node.Handle)
	*/

	return 0
}
