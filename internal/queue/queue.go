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

func (q *Queue[T]) Push(e T, idx int) {
	item := QueueItem[T]{
		item:  e,
		index: idx,
	}

	q.items = append(q.items, item)
}

func (q *Queue[T]) Pop() T {
	if q.IsEmpty() {
		var zero T
		return zero
	}
	e := q.items[0]
	q.items = q.items[1:]

	return e.item
}

func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}
