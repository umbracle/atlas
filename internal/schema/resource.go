package schema

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

type Resource struct {
	val cty.Value
}

func (r *Resource) getValue(s string) (cty.Value, bool) {
	parts := strings.Split(s, ".")

	path := cty.Path{}
	for _, part := range parts {
		num, err := strconv.Atoi(part)
		if err == nil {
			path = path.IndexInt(num)
		} else {
			path = path.GetAttr(part)
		}
	}

	val, err := path.Apply(r.val)
	if err != nil {
		return cty.NilVal, false
	}
	return val, true
}

func (r *Resource) Get(s string) string {
	val, _ := r.GetOk(s)
	return val
}

func (r *Resource) GetOk(s string) (string, bool) {
	val, ok := r.getValue(s)
	if !ok {
		return "", false
	}

	valStr, err := convert.Convert(val, cty.String)
	if err != nil {
		panic(fmt.Sprintf("BUG: cannot convert: %v", err))
	}
	return valStr.AsString(), true
}
