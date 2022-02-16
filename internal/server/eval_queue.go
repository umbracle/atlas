package server

import (
	"container/heap"
	"context"
	"sync"
	"time"

	"github.com/umbracle/atlas/internal/proto"
)

type evalQueue struct {
	lock     sync.Mutex
	heap     taskQueueImpl
	updateCh chan struct{}
}

func newEvalQueue() *evalQueue {
	return &evalQueue{
		heap:     taskQueueImpl{},
		updateCh: make(chan struct{}),
	}
}

func (e *evalQueue) add(eval *proto.Evaluation) {
	tt := &evalTask{
		eval:      eval,
		timestamp: time.Now(),
		ready:     true,
	}

	heap.Push(&e.heap, tt)

	select {
	case e.updateCh <- struct{}{}:
	default:
	}
}

func (e *evalQueue) popImpl() *proto.Evaluation {
	e.lock.Lock()
	if len(e.heap) != 0 {
		// pop the first value and remove it from the heap
		tt := heap.Pop(&e.heap).(*evalTask)
		e.lock.Unlock()

		return tt.eval
	}
	e.lock.Unlock()
	return nil
}

func (e *evalQueue) pop(ctx context.Context) *proto.Evaluation {
POP:
	tt := e.popImpl()
	if tt != nil {
		return tt
	}

	select {
	case <-e.updateCh:
		goto POP
	case <-ctx.Done():
		return nil
	}
}

type evalTask struct {
	eval      *proto.Evaluation
	index     int
	timestamp time.Time
	ready     bool
}

type taskQueueImpl []*evalTask

func (t taskQueueImpl) Len() int { return len(t) }

func (t taskQueueImpl) Less(i, j int) bool {
	iNoReady, jNoReady := !t[i].ready, !t[j].ready
	if iNoReady && jNoReady {
		return false
	} else if iNoReady {
		return false
	} else if jNoReady {
		return true
	}
	return t[i].timestamp.Before(t[j].timestamp)
}

func (t taskQueueImpl) Swap(i, j int) {
	t[i], t[j] = t[j], t[i]
	t[i].index = i
	t[j].index = j
}

func (t *taskQueueImpl) Push(x interface{}) {
	n := len(*t)
	item := x.(*evalTask)
	item.index = n
	*t = append(*t, item)
}

func (t *taskQueueImpl) Pop() interface{} {
	old := *t
	n := len(old)
	item := old[n-1]
	old[n-1] = nil
	item.index = -1
	*t = old[0 : n-1]
	return item
}
