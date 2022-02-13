package cmd

import (
	"os"

	"github.com/mitchellh/cli"
	"github.com/ryanuber/columnize"
	"github.com/umbracle/atlas/internal/proto"
	"google.golang.org/grpc"
)

// Commands returns the cli commands
func Commands() map[string]cli.CommandFactory {
	ui := &cli.BasicUi{
		Reader:      os.Stdin,
		Writer:      os.Stdout,
		ErrorWriter: os.Stderr,
	}

	meta := &Meta{
		UI: ui,
	}
	return map[string]cli.CommandFactory{
		"server": func() (cli.Command, error) {
			return &ServerCommand{
				UI: ui,
			}, nil
		},
		"deploy": func() (cli.Command, error) {
			return &DeployCommand{
				Meta: meta,
			}, nil
		},
		"version": func() (cli.Command, error) {
			return &VersionCommand{
				UI: ui,
			}, nil
		},
		"list": func() (cli.Command, error) {
			return &ListCommand{
				Meta: meta,
			}, nil
		},
		"stop": func() (cli.Command, error) {
			return &StopCommand{
				UI: ui,
			}, nil
		},
	}
}

func formatList(in []string) string {
	columnConf := columnize.DefaultConfig()
	columnConf.Empty = "<none>"
	return columnize.Format(in, columnConf)
}

func formatKV(in []string) string {
	columnConf := columnize.DefaultConfig()
	columnConf.Empty = "<none>"
	columnConf.Glue = " = "
	return columnize.Format(in, columnConf)
}

type Meta struct {
	UI cli.Ui
}

func (m *Meta) Conn() (proto.AtlasServiceClient, error) {
	conn, err := grpc.Dial("localhost:3030", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	return proto.NewAtlasServiceClient(conn), nil
}
