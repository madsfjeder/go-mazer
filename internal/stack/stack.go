// Package stack - Generic stack implementation with some helpers
package stack

import (
	"errors"
	"fmt"
	"slices"
)

type StackItem[T comparable] struct {
	Item  T
	Index int
}

type Stack[T comparable] struct {
	items []StackItem[T]
}

func New[T comparable]() *Stack[T] {
	return &Stack[T]{
		items: make([]StackItem[T], 0),
	}
}

func (s *Stack[T]) Push(e T, idx int) {
	item := StackItem[T]{
		e,
		idx,
	}
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, error) {
	if len(s.items) == 0 {
		var zero T
		return zero, errors.New("no more items")
	}
	lastIndex := len(s.items) - 1
	element := s.items[lastIndex]

	var zero StackItem[T]

	s.items[lastIndex] = zero
	s.items = s.items[:lastIndex]

	return element.Item, nil
}

func (s *Stack[T]) PopAll() []T {
	list := make([]T, 0)

	for _, v := range s.items {
		list = append(list, v.Item)
	}

	s.items = make([]StackItem[T], 0)

	return list
}

func (s *Stack[T]) PopAllWithIdx() []StackItem[T] {
	list := make([]StackItem[T], 0)

	for _, v := range s.items {
		list = append(list, StackItem[T]{
			Item:  v.Item,
			Index: v.Index,
		})
	}

	s.items = make([]StackItem[T], 0)

	return list
}

func (s Stack[T]) Print() {
	for _, v := range s.items {
		fmt.Println(v)
	}
}

func (s *Stack[T]) Copy() Stack[T] {
	itemsSlice := make([]StackItem[T], len(s.items))
	copy(itemsSlice, s.items)

	return Stack[T]{
		items: itemsSlice,
	}
}

// Filter in place
func (s *Stack[T]) Filter(fn func(e T) bool) {
	filteredItems := make([]StackItem[T], len(s.items))
	for _, v := range s.items {
		shouldAppend := fn(v.Item)
		if shouldAppend {
			filteredItems = append(filteredItems, v)
		}
	}

	s.items = filteredItems
}

func (s *Stack[T]) Reverse() {
	slices.Reverse(s.items)
}

func (s *Stack[T]) FindOrder(v T) int {
	for i, e := range s.items {
		if v == e.Item {
			return i
		}
	}

	return -1
}

func (s *Stack[T]) Length() int {
	return len(s.items)
}

func (s *Stack[T]) Items() []T {
	items := make([]T, 0)

	for _, v := range s.items {
		items = append(items, v.Item)
	}

	return items
}
