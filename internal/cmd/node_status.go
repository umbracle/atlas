package cmd

import (
	"context"
	"fmt"

	"github.com/umbracle/atlas/internal/proto"
)

// NodeStatusCommand is the command to show the version of the agent
type NodeStatusCommand struct {
	*Meta
}

// Help implements the cli.Command interface
func (c *NodeStatusCommand) Help() string {
	return `Usage: ensemble version
	
  Display the Ensemble version`
}

// Synopsis implements the cli.Command interface
func (c *NodeStatusCommand) Synopsis() string {
	return "Display the Ensemble version"
}

// Run implements the cli.Command interface
func (c *NodeStatusCommand) Run(args []string) int {

	id := args[0]

	client, err := c.Conn()
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	resp, err := client.NodeStatus(context.Background(), &proto.NodeStatusRequest{Id: id})
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	c.UI.Output(c.Colorize().Color(formatNodeStatus(resp)))
	return 0
}

func formatNodeStatus(r *proto.NodeStatusResponse) string {
	base := formatKV([]string{
		fmt.Sprintf("ID|%s", r.Node.Id),
		fmt.Sprintf("Chain|%s", r.Node.Chain),
		fmt.Sprintf("Running|%v", r.Node.Running),
	})

	if len(r.Events) != 0 {
		rows := make([]string, len(r.Events)+1)
		rows[0] = "Time|Message"
		for i, d := range r.Events {
			rows[i+1] = fmt.Sprintf("%s|%s",
				d.Timestamp.AsTime().String(),
				d.Message,
			)
		}
		base += "\n\n[bold]Events[reset]\n"
		base += formatList(rows)
	}
	return base
}
