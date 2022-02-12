package framework

import "github.com/umbracle/atlas/internal/proto"

type Plugin interface {
	Config() interface{}
	Chains() []string
	Build(i *Input) *proto.NodeSpec
}

type Input struct {
	Chain   string
	Datadir string
}
