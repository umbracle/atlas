package geth

import (
	"fmt"

	"github.com/umbracle/atlas/internal/framework"
	"github.com/umbracle/atlas/internal/proto"
	"github.com/umbracle/atlas/internal/schema"
)

const (
	imageName = "ethereum/client-go"
	imageTag  = "v1.9.25"
)

type Geth struct {
	config *config
}

type config struct {
	Cache int
}

func (g *Geth) Schema() *schema.Object {
	return &schema.Object{}
}

func (g *Geth) Config() interface{} {
	return &g.config
}

func (g *Geth) Chains() []string {
	return []string{
		"goerli",
	}
}

func (g *Geth) Build(input *framework.Input) *proto.NodeSpec {
	if g.config == nil {
		g.config = &config{}
	}

	fmt.Println(g.config)

	args := []string{
		"--datadir", "/data",
	}

	if input.Chain == "goerli" {
		args = append(args, "--goerli")
	} else {
		panic("bad")
	}

	return &proto.NodeSpec{
		Image: &proto.NodeSpec_Image{
			Image: imageName,
			Ref:   imageTag,
		},
		Args: args,
	}
}
