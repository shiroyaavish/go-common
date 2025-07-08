package linked_list

type LinkedList[Type any] struct {
	Total int
	Head  *Node[Type]
}

func (l *LinkedList[Type]) ConsumeCurrent() {
	current := l.Head
	l.Head = nil
	l.Head = current.Next
}

type Opts[Type any] interface {
	Set(l *LinkedList[Type])
}

func NewLinkedList[Type any](data []Type) *LinkedList[Type] {
	newList := new(LinkedList[Type])
	// Point head to first element
	for _, el := range data {
		if newList.Head == nil {
			newList.Head = newNode[Type](el)
		}
		current := newList.Head
		for current.Next != nil {
			current = current.Next
		}
		current.Next = newNode[Type](el)
	}
	data = nil
	return newList
}
