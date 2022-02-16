package proto

import "reflect"

func (n *NodeSpec) Equal(nn *NodeSpec) bool {
	if !reflect.DeepEqual(n.Args, nn.Args) {
		return false
	}
	if !reflect.DeepEqual(n.Image, nn.Image) {
		return false
	}
	return true
}
