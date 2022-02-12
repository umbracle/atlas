package plugins

import (
	"github.com/umbracle/atlas/internal/framework"
	"github.com/umbracle/atlas/plugins/geth"
)

var Plugins = map[string]framework.Plugin{
	"geth": &geth.Geth{},
}
