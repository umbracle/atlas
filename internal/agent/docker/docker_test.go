package docker

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/umbracle/atlas/internal/proto"
)

func TestDocker_ListContainers(t *testing.T) {
	d, err := NewDocker()
	assert.NoError(t, err)

	spec := &proto.NodeSpec{
		Image: &proto.NodeSpec_Image{
			Image: "ethereum/client-go",
			Ref:   "v1.9.25",
		},
		Args: []string{
			"--dev",
		},
	}
	respID, err := d.Run(context.Background(), spec)
	assert.NoError(t, err)

	defer d.StopID(respID)

	containers, err := d.ListContainers()
	assert.NoError(t, err)

	assert.True(t, spec.Equal(containers[0].Spec))
}
