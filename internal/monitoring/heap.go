package monitoring

import "container/heap"

type Heap []*TimedURL

// should not be called directly
func (h *Heap) Len() int {
	return len(*h)
}

// should not be called directly
func (h *Heap) Less(i, j int) bool {
	return (*h)[i].callTime.Before((*h)[j].callTime)
}

// should not be called directly
func (h *Heap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]

	(*h)[i].index = i
	(*h)[j].index = j
}

// should not be called directly
func (h *Heap) Push(x any) {
	item := x.(*TimedURL)
	item.index = len(*h)
	*h = append(*h, item)
}

// should not be called directly
func (h *Heap) Pop() any {
	old := *h
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // avoid memory leak
	item.index = -1 // for safety
	*h = old[0 : n-1]

	// free unused memory
	if len(*h) < cap(*h)/2 {
		*h = append(Heap(nil), *h...)
	}

	return item
}

// should not be called directly
func (h *Heap) Peek() *TimedURL {
	return (*h)[0]
}

func NewHeap(urls ...*TimedURL) *Heap {
	h := Heap(urls)
	heap.Init(&h)
	return &h
}
