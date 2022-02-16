package cmd

import (
	"context"
	"fmt"

	"github.com/umbracle/atlas/internal/proto"
)

// ProvidersListCommand is the command to show the version of the agent
type ProvidersListCommand struct {
	*Meta
}

// Help implements the cli.Command interface
func (c *ProvidersListCommand) Help() string {
	return `Usage: ensemble version
	
  Display the Ensemble version`
}

// Synopsis implements the cli.Command interface
func (c *ProvidersListCommand) Synopsis() string {
	return "Display the Ensemble version"
}

// Run implements the cli.Command interface
func (c *ProvidersListCommand) Run(args []string) int {

	conn, err := c.Conn()
	if err != nil {
		panic(err)
	}

	resp, err := conn.ListProviders(context.Background(), &proto.ListProvidersRequest{})
	if err != nil {
		panic(err)
	}
	c.UI.Output(formatProviders(resp.Providers))
	return 0
}

func formatProviders(deps []*proto.Provider) string {
	if len(deps) == 0 {
		return "No providers found"
	}

	rows := make([]string, len(deps)+1)
	rows[0] = "Id|Name|Provider"
	for i, d := range deps {
		rows[i+1] = fmt.Sprintf("%s|%s|%s",
			d.Id,
			d.Name,
			d.Provider,
		)
	}
	return formatList(rows)
}
