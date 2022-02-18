package framework

import (
	"github.com/umbracle/atlas/internal/proto"
	"github.com/umbracle/atlas/internal/schema"
)

type Plugin interface {
	Config() interface{}
	Schema() *schema.Object
	Chains() []string
	Build(i *Input) *proto.NodeSpec
}

type Input struct {
	Chain   string
	Datadir string
}
