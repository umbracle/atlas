package schema

import (
	"fmt"
	"testing"
)

func TestOpenAPIV3Marshal(t *testing.T) {
	cases := []struct {
		filename string
		schema   *Object
	}{
		{
			"openapiv3.json",
			&Object{
				Fields: map[string]*Field{
					"a": {
						Type: &String,
					},
					"b": {
						Type: &Object{
							Fields: map[string]*Field{
								"c": {
									Type:     &Number,
									Required: true,
								},
							},
						},
					},
					"d": {
						Type: &String,
					},
					"e": {
						Type: &Array{
							Elem: &Object{
								Fields: map[string]*Field{
									"e1": {
										Type: &String,
									},
									"e2": {
										Type: &String,
									},
								},
							},
						},
					},
				},
			},
		},
	}
	for _, c := range cases {
		res, err := c.schema.OpenAPIV3JSON()
		if err != nil {
			t.Fatal(err)
		}
		fmt.Println(string(res))
	}
}
