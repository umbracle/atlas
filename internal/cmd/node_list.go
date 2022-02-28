package cmd

import (
	"context"
	"fmt"

	"github.com/umbracle/atlas/internal/proto"
)

// NodeListCommand is the command to show the version of the agent
type NodeListCommand struct {
	*Meta
}

// Help implements the cli.Command interface
func (c *NodeListCommand) Help() string {
	return `Usage: ensemble version
	
  Display the Ensemble version`
}

// Synopsis implements the cli.Command interface
func (c *NodeListCommand) Synopsis() string {
	return "Display the Ensemble version"
}

// Run implements the cli.Command interface
func (c *NodeListCommand) Run(args []string) int {
	client, err := c.Conn()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	resp, err := client.ListNodes(context.Background(), &proto.ListNodesRequest{})
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output(formatNodes(resp.Nodes))
	return 0
}

func formatNodes(deps []*proto.Node) string {
	if len(deps) == 0 {
		return "No nodes found"
	}

	rows := make([]string, len(deps)+1)
	rows[0] = "Name|Chain|Running"
	for i, d := range deps {
		rows[i+1] = fmt.Sprintf("%s|%s|%v",
			d.Id,
			d.Chain,
			d.Running,
		)
	}
	return formatList(rows)
}
