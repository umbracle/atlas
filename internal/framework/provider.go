package framework

import (
	"context"

	"github.com/umbracle/atlas/internal/proto"
	"github.com/umbracle/atlas/internal/schema"
)

type Provider interface {
	Init()
	Config() (interface{}, error)
	Update(ctx context.Context, node *proto.Node) error
	Schema() *schema.Object
}
