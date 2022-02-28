package framework

import (
	"context"

	"github.com/umbracle/atlas/internal/proto"
	"github.com/umbracle/atlas/internal/schema"
)

type Provider interface {
	Init()
	Config() interface{}
	Update(ctx context.Context, old, new interface{}, node *proto.Node) error
	Schema() *schema.Object
}

type Framework struct {
}
