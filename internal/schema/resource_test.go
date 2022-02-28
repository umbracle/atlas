package schema

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResource(t *testing.T) {

	s := &Object{
		Fields: map[string]*Field{
			"a": {
				Type: &Object{
					Fields: map[string]*Field{
						"d": {
							Type: &Number,
						},
					},
				},
			},
			"b": {
				Type: &Array{
					Elem: &Number,
				},
			},
		},
	}

	oldVal := map[string]interface{}{
		"a": map[string]int{
			"d": 2,
		},
		"b": []int{
			3,
		},
	}
	curVal := map[string]interface{}{
		"a": map[string]int{
			"d": 2,
		},
		"b": []int{
			4,
		},
	}

	res, err := NewResource(s, oldVal, curVal)
	assert.NoError(t, err)

	fmt.Println(res)

	fmt.Println(res.HasChanged("a.b"))
}
