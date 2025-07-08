package linked_list

type Node[Type any] struct {
	Next *Node[Type]
	data *Type
}

func (n *Node[Type]) Get() *Type {
	return n.data
}

func newNode[Type any](data Type) *Node[Type] {
	return &Node[Type]{
		data: &data,
		Next: nil,
	}
}
