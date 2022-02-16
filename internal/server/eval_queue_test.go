package server

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/umbracle/atlas/internal/proto"
)

func TestEvalQueue_Simple(t *testing.T) {
	eval := newEvalQueue()
	eval.add(&proto.Evaluation{
		Node: "a",
	})
	eval.add(&proto.Evaluation{
		Node: "b",
	})

	assert.Equal(t, "a", eval.popImpl().Node)
	assert.Equal(t, "b", eval.popImpl().Node)
	assert.Nil(t, eval.popImpl())
}

func TestEvalQueue_Pop(t *testing.T) {
	eval := newEvalQueue()

	evalCh := make(chan *proto.Evaluation)
	go func() {
		for {
			val := eval.pop(context.Background())
			evalCh <- val
		}
	}()

	eval.add(&proto.Evaluation{
		Node: "a",
	})

	select {
	case val := <-evalCh:
		assert.Equal(t, "a", val.Node)
	case <-time.After(1 * time.Second):
		t.Fatal("timeout")
	}
}
