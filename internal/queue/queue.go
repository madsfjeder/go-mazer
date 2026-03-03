// Package queue implementation of basic queue structure
package queue

type QueueItem[T comparable] struct {
	item  T
	index int
}

type Queue[T comparable] struct {
	items []QueueItem[T]
}

func New[T comparable]() *Queue[T] {
	return &Queue[T]{
		items: make([]QueueItem[T], 0),
	}
}
