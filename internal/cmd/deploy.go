package cmd

import (
	"context"
	"flag"
	"fmt"

	"github.com/umbracle/atlas/internal/proto"
)

// DeployCommand is the command to show the version of the agent
type DeployCommand struct {
	*Meta

	chain    string
	config   string
	provider string
	args     string
	id       string
}

// Help implements the cli.Command interface
func (c *DeployCommand) Help() string {
	return `Usage: ensemble version
	
  Display the Ensemble version`
}

// Synopsis implements the cli.Command interface
func (c *DeployCommand) Synopsis() string {
	return "Display the Ensemble version"
}

// Run implements the cli.Command interface
func (c *DeployCommand) Run(args []string) int {
	fmt.Println("- agent -")

	flags := flag.NewFlagSet("deploy", flag.PanicOnError)
	flags.StringVar(&c.chain, "chain", "", "")
	flags.StringVar(&c.config, "config", "", "")
	flags.StringVar(&c.provider, "provider", "", "")
	flags.StringVar(&c.args, "args", "{}", "") // these are the provider args
	flags.StringVar(&c.id, "id", "", "")

	if err := flags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	conn, err := c.Conn()
	if err != nil {
		panic(err)
	}

	provider, err := conn.GetProviderByName(context.Background(), &proto.GetProviderByNameRequest{Name: c.provider})
	if err != nil {
		panic(err)
	}

	req := &proto.DeployRequest{
		Chain:      c.chain,
		Config:     c.config,
		Plugin:     "geth",
		ProviderId: provider.Id,
		Args:       c.args,
		Id:         c.id,
	}
	resp, err := conn.Deploy(context.Background(), req)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	fmt.Println(resp.Node)
	return 0
}
