package schema

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/zclconf/go-cty/cty"
	"gopkg.in/yaml.v2"
)

// OpenAPIV3JSON creates the OpenAPIv3 spec for the Schema in JSON format
func (s *Object) OpenAPIV3JSON() ([]byte, error) {
	return json.Marshal(s.OpenAPIV3())
}

// OpenAPIV3Yaml creates the OpenAPIv3 spec for the Schema in YAML format
func (s *Object) OpenAPIV3Yaml() ([]byte, error) {
	return yaml.Marshal(s.OpenAPIV3())
}

// OpenAPIV3 returns the OpenApi v3 representation of the schema
func (s *Object) OpenAPIV3() interface{} {
	return marshalOpenAPIRecord(s)
}

func marshalOpenAPIRecord(r *Object) map[string]interface{} {
	res := map[string]interface{}{
		"type": "object",
	}

	properties := map[string]interface{}{}
	required := []string{}

	for k, v := range r.Fields {
		aux := marshalOpenAPIType(v.Type)

		if v.Required {
			required = append(required, k)
		}
		properties[k] = aux
	}

	if len(required) != 0 {
		res["required"] = required
	}
	res["properties"] = properties
	return res
}

func typeToOpenAPIType(t Type) string {
	switch obj := t.(type) {
	case *Object:
		return "object"

	case *Array:
		return "array"

	case *basicType:
		switch obj.Type() {
		case cty.Number:
			return "integer"

		case cty.Bool:
			return "boolean"

		case cty.String:
			return "string"
		}
	}
	panic(fmt.Sprintf("BUG: type not found: %s", reflect.TypeOf(t)))
}

func marshalOpenAPIType(t Type) map[string]interface{} {
	res := map[string]interface{}{}

	// add the type
	res["type"] = typeToOpenAPIType(t)

	switch obj := t.(type) {
	case *Object:
		// print object
		return marshalOpenAPIRecord(obj)

	case *Array:
		// print array
		res["items"] = marshalOpenAPIType(obj.Elem)
	}

	return res
}
