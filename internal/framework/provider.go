package framework

import (
	"context"

	"github.com/umbracle/atlas/internal/proto"
)

type Provider interface {
	Config() (interface{}, error)
	Update(ctx context.Context, node *proto.Node) error
}
