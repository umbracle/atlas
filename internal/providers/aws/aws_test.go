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
		Id:             "a",
		ExpectedConfig: `{"type": "t2.small"}`,
	}
	assert.NoError(t, a.Update(context.Background(), node))

	// Update the node
	node.ExpectedConfig = `{"type": "t2.medium"}`
	assert.NoError(t, a.Update(context.Background(), node))
}
