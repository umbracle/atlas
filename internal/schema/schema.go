package schema

import (
	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/gocty"
)

type Type interface {
	Type() cty.Type
}

type Object struct {
	Fields map[string]*Field
}

func (o *Object) DecodeValue(val interface{}) (cty.Value, error) {
	return gocty.ToCtyValue(val, o.Type())
}

func (o *Object) Type() cty.Type {
	attrTypes := map[string]cty.Type{}
	for name, field := range o.Fields {
		attrTypes[name] = field.Type.Type()
	}
	return cty.Object(attrTypes)
}

type Field struct {
	Type        Type
	Description string
	Required    bool
	Default     interface{}
}

type Array struct {
	Elem Type
}

func (a *Array) Type() cty.Type {
	return cty.List(a.Elem.Type())
}

type basicType cty.Type

var (
	String basicType = basicType(cty.String)
	Number basicType = basicType(cty.Number)
	Bool   basicType = basicType(cty.Bool)
)

func (b *basicType) Type() cty.Type {
	return cty.Type(*b)
}
