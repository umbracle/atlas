package runtime

import (
	"context"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/volume"
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

func (d *Docker) Run(ctx context.Context, spec *proto.NodeSpec) (*proto.Node_Handle, error) {

	// create random volume
	d.cli.VolumeCreate(ctx, volume.VolumeCreateBody{})

	// create container
	resp, err := d.cli.ContainerCreate(ctx, &container.Config{
		Image: spec.Image.Image + ":" + spec.Image.Ref,
		Cmd:   spec.Args,
	}, nil, nil, nil, "")
	if err != nil {
		return nil, err
	}

	if err := d.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return nil, err
	}
	handle := &proto.Node_Handle{
		Handle: resp.ID,
	}
	return handle, nil
}

func (d *Docker) Stop(spec *proto.Node_Handle) error {
	return d.cli.ContainerStop(context.Background(), spec.Handle, nil)
}
