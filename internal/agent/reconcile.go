package agent

import (
	"github.com/umbracle/atlas/internal/agent/docker"
	"github.com/umbracle/atlas/internal/proto"
)

type reconciler struct {
	expected *proto.NodeSpec
	found    *docker.RunningContainer
}

type plan struct {
	stop  *string
	start *proto.NodeSpec
}

func (p *plan) empty() bool {
	return p.stop == nil && p.start == nil
}

func (r *reconciler) reconcile() *plan {
	plan := &plan{}

	// assume expect is never empty
	if r.expected.Expected == proto.NodeSpec_Running {
		if r.found == nil {
			plan.start = r.expected
		}
	} else if r.expected.Expected == proto.NodeSpec_Terminated {
		if r.found != nil {
			plan.stop = &r.found.Id
		}
	}

	return plan
}
