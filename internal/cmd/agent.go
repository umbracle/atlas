package cmd

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/mitchellh/cli"
	"github.com/umbracle/atlas/internal/agent"
)

// AgentCommand is the command to show the version of the agent
type AgentCommand struct {
	UI cli.Ui
}

// Help implements the cli.Command interface
func (c *AgentCommand) Help() string {
	return `Usage: ensemble version
	
  Display the Ensemble version`
}

// Synopsis implements the cli.Command interface
func (c *AgentCommand) Synopsis() string {
	return "Display the Ensemble version"
}

// Run implements the cli.Command interface
func (c *AgentCommand) Run(args []string) int {

	/*
		readKey, err := ioutil.ReadFile("/home/ferran/Downloads/atlas.pem")
		if err != nil {
			panic(err)
		}
		key, err := ssh.ParsePrivateKey(readKey)
		if err != nil {
			panic(err)
		}

		// Authentication
		config := &ssh.ClientConfig{
			User: "ec2-user",
			// https://github.com/golang/go/issues/19767
			// as clientConfig is non-permissive by default
			// you can set ssh.InsercureIgnoreHostKey to allow any host
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(key),
			},
			//alternatively, you could use a password
			Timeout: 10 * time.Second,
		}
		// Connect
		client, err := ssh.Dial("tcp", net.JoinHostPort("34.216.68.253", "22"), config)
		if err != nil {
			panic(err)
		}
		// Create a session. It is one session per command.
		session, err := client.NewSession()
		if err != nil {
			panic(err)
		}

		fmt.Println(session)

		session.Stdout = os.Stdout
		session.Stderr = os.Stdout
		session.Run("docker ps")

		return 0
	*/

	logger := hclog.New(&hclog.LoggerOptions{
		Name:  "agent",
		Level: hclog.LevelFromString("info"),
	})

	config := &agent.Config{}

	agent, err := agent.NewAgent(logger, config)
	if err != nil {
		c.UI.Error(err.Error())
		return 1
	}

	logger.Info("Agent running")
	return c.handleSignals(agent.Close)
}

func (c *AgentCommand) handleSignals(closeFn func()) int {
	signalCh := make(chan os.Signal, 4)
	signal.Notify(signalCh, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)

	sig := <-signalCh

	c.UI.Output(fmt.Sprintf("Caught signal: %v", sig))
	c.UI.Output("Gracefully shutting down agent...")

	gracefulCh := make(chan struct{})
	go func() {
		if closeFn != nil {
			closeFn()
		}
		close(gracefulCh)
	}()

	select {
	case <-signalCh:
		return 1
	case <-time.After(5 * time.Second):
		return 1
	case <-gracefulCh:
		return 0
	}
}
