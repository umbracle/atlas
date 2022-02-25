package docker

import (
	"context"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/container"
	"github.com/docker/docker/api/types/filters"
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

func (d *Docker) PullImage(ctx context.Context, spec *proto.NodeSpec) (bool, error) {
	canonicalName := "docker.io/" + spec.Image.Image + ":" + spec.Image.Ref

	fmt.Println("- pull -")
	fmt.Println(canonicalName)

	_, _, err := d.cli.ImageInspectWithRaw(ctx, canonicalName)
	if err != nil {
		fmt.Println("- not found -")

		reader, err := d.cli.ImagePull(ctx, canonicalName, types.ImagePullOptions{})
		if err != nil {
			return false, err
		}
		_, err = io.Copy(os.Stdout, reader)
		if err != nil {
			return false, err
		}
		return true, nil
	}
	return false, nil
}

func (d *Docker) Run(ctx context.Context, spec *proto.NodeSpec) (string, error) {
	// create container
	config := &container.Config{
		Image: spec.Image.Image + ":" + spec.Image.Ref,
		Cmd:   spec.Args,
		Labels: map[string]string{
			"atlas": "true",
		},
	}
	hostConfig := &container.HostConfig{
		Binds: []string{
			"/data:/data",
		},
	}
	resp, err := d.cli.ContainerCreate(ctx, config, hostConfig, nil, nil, "")
	if err != nil {
		return "", err
	}

	if err := d.cli.ContainerStart(ctx, resp.ID, types.ContainerStartOptions{}); err != nil {
		return "", err
	}
	fmt.Println(resp.ID)
	return resp.ID, nil
}

type RunningContainer struct {
	Id   string
	Spec *proto.NodeSpec
}

func (d *Docker) ListContainers() ([]*RunningContainer, error) {
	filters := filters.NewArgs()
	filters.Add("label", "atlas")

	containers, err := d.cli.ContainerList(context.Background(), types.ContainerListOptions{
		Filters: filters,
	})

	res := []*RunningContainer{}
	for _, c := range containers {

		data, err := d.cli.ContainerInspect(context.Background(), c.ID)
		if err != nil {
			// this should not happen right?
			panic(err)
		}

		prts := strings.Split(c.Image, ":")

		// convert container into a proto.NodeSpec
		spec := &proto.NodeSpec{
			Image: &proto.NodeSpec_Image{
				Image: prts[0],
				Ref:   prts[1],
			},
			Args: data.Args,
		}

		res = append(res, &RunningContainer{
			Id:   c.ID,
			Spec: spec,
		})
	}
	return res, err
}

func (d *Docker) StopID(id string) error {
	return d.cli.ContainerRemove(context.Background(), id, types.ContainerRemoveOptions{Force: true})
}

func (d *Docker) Stop(spec *proto.Node_Handle) error {
	return d.cli.ContainerStop(context.Background(), spec.Handle, nil)
}
