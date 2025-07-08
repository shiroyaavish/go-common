package stack

import "sync"

type Stack[Item any] struct {
	items []Item
	lock  sync.RWMutex
}

func New[Item any]() Stack[Item] {
	return Stack[Item]{
		items: make([]Item, 0),
		lock:  sync.RWMutex{},
	}
}

// Push adds an Item to the top of the stack
func (s *Stack[Item]) Push(i Item) {
	s.lock.Lock()
	s.items = append(s.items, i)
	s.lock.Unlock()
}

// Pop removes an Item from the top of the stack
func (s *Stack[Item]) Pop() *Item {
	s.lock.Lock()
	item := new(Item)
	if len(s.items) == 0 {
		s.lock.Unlock()
		return nil
	}

	*item, s.items = s.items[len(s.items)-1], s.items[:len(s.items)-1]
	s.lock.Unlock()
	return item
}
