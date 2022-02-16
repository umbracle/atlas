package agent

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/umbracle/atlas/internal/agent/docker"
	"github.com/umbracle/atlas/internal/proto"
)

func TestReconcile_Nil_Running(t *testing.T) {
	r := &reconciler{
		expected: &proto.NodeSpec{
			Args:     []string{"a"},
			Expected: proto.NodeSpec_Running,
		},
		found: &docker.RunningContainer{
			Id: "b",
		},
	}

	p := r.reconcile()
	assert.True(t, p.empty())
}

func TestReconcile_Nil_Terminated(t *testing.T) {
	r := &reconciler{
		expected: &proto.NodeSpec{
			Args:     []string{"a"},
			Expected: proto.NodeSpec_Terminated,
		},
	}

	p := r.reconcile()
	assert.True(t, p.empty())
}

func TestReconcile_Create(t *testing.T) {
	r := &reconciler{
		expected: &proto.NodeSpec{
			Args:     []string{"a"},
			Expected: proto.NodeSpec_Running,
		},
	}

	p := r.reconcile()
	assert.NotNil(t, p.start)
	assert.Nil(t, p.stop)
}

func TestReconcile_Terminate(t *testing.T) {
	r := &reconciler{
		expected: &proto.NodeSpec{
			Args:     []string{"a"},
			Expected: proto.NodeSpec_Terminated,
		},
		found: &docker.RunningContainer{
			Id: "b",
		},
	}

	p := r.reconcile()
	assert.Nil(t, p.start)
	assert.NotNil(t, p.stop)
}
