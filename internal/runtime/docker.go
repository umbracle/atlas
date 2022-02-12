package runtime

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/client"
	"github.com/umbracle/atlas/internal/proto"
)

// Docker is a sugarcoat version of the docker client
type Docker struct {
	cli *client.Client
}

func NewDocker() (*Docker, error) {
	cli, err := client.NewClientWithOpts(client.FromEnv, client.WithAPIVersionNegotiation())
	if err != nil {
		return nil, err
	}
	client := &Docker{
		cli: cli,
	}
	return client, nil
}

func (d *Docker) Run(ctx context.Context, spec *proto.NodeSpec) {
	resp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image: spec.Image.Image + ":" + spec.Image.Ref,
		Cmd:   spec.Args,
	}, nil, nil, nil, "")
	if err != nil {
		panic(err)
	}

	if err := d.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		panic(err)
	}
}
