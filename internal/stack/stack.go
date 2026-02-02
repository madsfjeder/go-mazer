// Package stack - Generic stack implementation with some helpers
package stack

import (
	"errors"
	"fmt"
	"slices"
)

type Stack[T comparable] struct {
	items []T
}

func New[T comparable]() *Stack[T] {
	return &Stack[T]{
		items: make([]T, 0),
	}
}

func (s *Stack[T]) Push(e T) {
	s.items = append(s.items, e)
}

func (s *Stack[T]) Pop() (T, error) {
	if len(s.items) == 0 {
		var zero T
		return zero, errors.New("no more items")
	}
	lastIndex := len(s.items) - 1
	element := s.items[lastIndex]

	var zero T

	s.items[lastIndex] = zero
	s.items = s.items[:lastIndex]

	return element, nil
}

func (s Stack[T]) Print() {
	for _, v := range s.items {
		fmt.Println(v)
	}
}

func (s Stack[T]) Copy() Stack[T] {
	itemsSlice := make([]T, len(s.items))
	copy(itemsSlice, s.items)

	return Stack[T]{
		items: itemsSlice,
	}
}

func (s *Stack[T]) Reverse() {
	slices.Reverse(s.items)
}

func (s *Stack[T]) FindOrder(v T) int {
	for i, e := range s.items {
		if v == e {
			return i
		}
	}

	return -1
}

func (s *Stack[T]) Length() int {
	return len(s.items)
}
