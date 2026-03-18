// Package pqueue - Priority queue implementation (for AStar)
package pqueue

import "fmt"

type PQueueItem[T comparable] struct {
	item  T
	value float32
}

type PQueue[T comparable] struct {
	items []PQueueItem[T]
}

func New[T comparable]() *PQueue[T] {
	return &PQueue[T]{
		items: make([]PQueueItem[T], 0),
	}
}

func (p *PQueue[T]) Insert(e T, v float32) {
	item := PQueueItem[T]{
		item:  e,
		value: v,
	}
	p.items = append(p.items, item)
	p.shiftUp(len(p.items) - 1)
}

func (p *PQueue[T]) Pop() T {
	size := len(p.items)

	if size == 0 {
		var zero PQueueItem[T]

		return zero.item
	}

	result := p.items[0]
	p.items[0] = p.items[size-1]
	p.items = p.items[0 : len(p.items)-1]
	p.shiftDown(0, len(p.items))
	return result.item
}

func (p *PQueue[T]) PopWithValue() PQueueItem[T] {
	size := len(p.items)

	if size == 0 {
		var zero PQueueItem[T]

		return zero
	}

	result := p.items[0]
	p.items[0] = p.items[size-1]
	p.items = p.items[0 : len(p.items)-1]
	p.shiftDown(0, len(p.items))
	return result
}

func (p PQueue[T]) Length() int {
	return len(p.items)
}

func (p PQueue[T]) PrintAll() {
	for i := range p.items {
		item := p.PopWithValue()
		fmt.Println("Item ", i, item.value)
	}
}

func (p *PQueue[T]) shiftUp(i int) {
	for i > 0 && p.items[parent(i)].value > p.items[i].value {
		p.items[parent(i)], p.items[i] = p.items[i], p.items[parent(i)]
		i = parent(i)
	}
}

func (p *PQueue[T]) shiftDown(i, size int) {
	maxIndex := i
	l := p.leftChild(i)

	if l < size && p.items[l].value < p.items[maxIndex].value {
		maxIndex = l
	}

	r := p.rightChild(i)

	if r < size && p.items[r].value < p.items[maxIndex].value {
		maxIndex = r
	}

	if i != maxIndex {
		p.items[i], p.items[maxIndex] = p.items[maxIndex], p.items[i]
		p.shiftDown(maxIndex, size)
	}
}

func (p *PQueue[T]) leftChild(i int) int {
	return 2*i + 1
}

func (p *PQueue[T]) rightChild(i int) int {
	return 2*i + 2
}

func parent(i int) int {
	return (i - 1) / 2
}
