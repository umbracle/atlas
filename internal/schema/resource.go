package schema

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/zclconf/go-cty/cty"
	"github.com/zclconf/go-cty/cty/convert"
)

type Resource struct {
	sch *Object
	old cty.Value
	val cty.Value
}

func NewResource(sch *Object, old, cur map[string]interface{}) (*Resource, error) {
	oldVal, err := sch.DecodeValue(old)
	if err != nil {
		return nil, err
	}
	curVal, err := sch.DecodeValue(cur)
	if err != nil {
		return nil, err
	}

	resource := &Resource{
		sch: sch,
		old: oldVal,
		val: curVal,
	}
	return resource, nil
}

func (r *Resource) getValue(refVal cty.Value, s string) (cty.Value, bool) {
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

	val, err := path.Apply(refVal)
	if err != nil {
		return cty.NilVal, false
	}
	return val, true
}

func (r *Resource) HasChanged(s string) bool {
	val0, ok0 := r.getValue(r.val, s)
	val1, ok1 := r.getValue(r.old, s)

	if ok0 != ok1 {
		return true
	}
	return !val0.RawEquals(val1)
}

func (r *Resource) Get(s string) string {
	val, _ := r.GetOk(s)
	return val
}

func (r *Resource) GetOk(s string) (string, bool) {
	val, ok := r.getValue(r.val, s)
	if !ok {
		return "", false
	}

	valStr, err := convert.Convert(val, cty.String)
	if err != nil {
		panic(fmt.Sprintf("BUG: cannot convert: %v", err))
	}
	return valStr.AsString(), true
}
