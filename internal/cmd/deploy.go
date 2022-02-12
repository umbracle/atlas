package cmd

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"

	"github.com/google/uuid"
	"github.com/mitchellh/cli"
	"github.com/umbracle/atlas/internal/framework"
	"github.com/umbracle/atlas/internal/proto"
	"github.com/umbracle/atlas/internal/runtime"
	"github.com/umbracle/atlas/internal/state"
	"github.com/umbracle/atlas/plugins"
)

// DeployCommand is the command to show the version of the agent
type DeployCommand struct {
	UI cli.Ui

	chain  string
	config string
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

	if err := flags.Parse(args); err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	d, err := runtime.NewDocker()
	if err != nil {
		panic(err)
	}

	plugin, ok := plugins.Plugins["geth"]
	if !ok {
		panic("plugin not found")

	}

	state, err := state.NewState("./state.db")
	if err != nil {
		panic(err)
	}

	// check if chain exists
	existsChain := false
	for _, chain := range plugin.Chains() {
		if chain == c.chain {
			existsChain = true
		}
	}
	if !existsChain {
		c.UI.Error(fmt.Sprintf("chain %s is not available", c.chain))
		return 1
	}

	if c.config != "" {
		configRaw, err := ioutil.ReadFile("./config.json")
		if err != nil {
			panic(err)
		}

		config := plugin.Config()
		if err := json.Unmarshal(configRaw, &config); err != nil {
			panic(err)
		}
	}

	input := &framework.Input{
		Chain: c.chain,
	}
	nodeSpec := plugin.Build(input)

	id := UUID()
	node := &proto.Node{
		Id:   id,
		Spec: nodeSpec,
	}

	if err := state.UpsertNode(node); err != nil {
		panic(err)
	}
	d.Run(context.Background(), nodeSpec)

	return 0
}

func UUID() string {
	return uuid.New().String()
}
