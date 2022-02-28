package server

import "testing"

func TestScheduler(t *testing.T) {
	sched := &scheduler{}
	sched.process()
}
