package schema

import (
	"fmt"
	"testing"
)

func TestSchema(t *testing.T) {

	s := &Object{
		Fields: map[string]*Field{
			"a": {
				Type: &Array{
					Elem: &Object{
						Fields: map[string]*Field{
							"d": {
								Type: &Number,
							},
						},
					},
				},
			},
			"b": {
				Type: &Object{
					Fields: map[string]*Field{
						"c": {
							Type: &Number,
						},
					},
				},
			},
		},
	}

	data, _ := s.Type().MarshalJSON()
	fmt.Println(string(data))

	rawVal := map[string]interface{}{
		"a": []interface{}{
			map[string]int{
				"d": 1,
			},
			map[string]int{
				"d": 2,
			},
		},
		"b": map[string]int{
			"c": 3,
		},
	}
	val, err := s.DecodeValue(rawVal)
	if err != nil {
		t.Fatal(err)
	}

	r := &Resource{
		val: val,
	}
	fmt.Println(r.Get("b.c"))
}
