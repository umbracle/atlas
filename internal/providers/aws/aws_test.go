package aws

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/umbracle/atlas/internal/proto"
)

func TestAws(t *testing.T) {
	a := &AwsProvider{}
	a.Init()

	// Create the node
	node := &proto.Node{
		Id: "a",
	}
	config0 := &config{
		Type: "t2.small",
	}
	assert.NoError(t, a.Update(context.Background(), nil, config0, node))

	// Update the node
	config1 := &config{
		Type: "t2.medium",
	}
	assert.NoError(t, a.Update(context.Background(), config0, config1, node))
}
