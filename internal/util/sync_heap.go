package util

import (
	"container/heap"
	"sync"
)

type CustomHeapInterface[T any] interface {
	heap.Interface
	Peek() T
}

type SyncHeap[T any] struct {
	h     CustomHeapInterface[T]
	mutex sync.Mutex
}

func (sh *SyncHeap[T]) Peek() T {
	sh.mutex.Lock()
	defer sh.mutex.Unlock()

	return sh.h.Peek()
}

func (sh *SyncHeap[T]) Fix(i int) {
	sh.mutex.Lock()
	defer sh.mutex.Unlock()

	heap.Fix(sh.h, i)
}

func (sh *SyncHeap[T]) Push(x T) {
	sh.mutex.Lock()
	defer sh.mutex.Unlock()

	heap.Push(sh.h, x)
}

func (sh *SyncHeap[T]) Pop() T {
	sh.mutex.Lock()
	defer sh.mutex.Unlock()

	return heap.Pop(sh.h).(T)
}

func (sh *SyncHeap[T]) Len() int {
	sh.mutex.Lock()
	defer sh.mutex.Unlock()

	return sh.h.Len()
}

func NewSyncHeap[T any](h CustomHeapInterface[T]) *SyncHeap[T] {
	return &SyncHeap[T]{
		h: h,
	}
}
