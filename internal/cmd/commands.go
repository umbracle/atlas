package cmd

import (
	"os"

	"github.com/mitchellh/cli"
)

// Commands returns the cli commands
func Commands() map[string]cli.CommandFactory {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	return map[string]cli.CommandFactory{
		"deploy": func() (cli.Command, error) {
			return &DeployCommand{
				UI: ui,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &VersionCommand{
				UI: ui,
			}, nil
		},
		"list": func() (cli.Command, error) {
			return &ListCommand{
				UI: ui,
			}, nil
		},
	}
}
