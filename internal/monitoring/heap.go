package monitoring

import (
	"container/heap"
	"time"
)

type TimedURL struct {
	URL      string
	UserId   uint32
	Interval time.Duration
	callTime time.Time
	index    int
}

func NewTimeURL(URL string, UserId uint32, Interval time.Duration) *TimedURL {
	return &TimedURL{
		URL:      URL,
		UserId:   UserId,
		Interval: Interval,
	}
}

type Heap []*TimedURL

func (h *Heap) Len() int {
	return len(*h)
}

func (h *Heap) Less(i, j int) bool {
	return (*h)[i].callTime.Before((*h)[j].callTime)
}

func (h *Heap) Swap(i, j int) {
	(*h)[i], (*h)[j] = (*h)[j], (*h)[i]

	(*h)[i].index = i
	(*h)[j].index = j
}

// should not be call directly
func (h *Heap) Push(x any) {
	item := x.(*TimedURL)
	item.index = len(*h)
	*h = append(*h, item)
}

// should not be call directly
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

func (h *Heap) Peek() *TimedURL {
	return (*h)[0]
}

func (h *Heap) PopRoot() *TimedURL {
	return heap.Pop(h).(*TimedURL)
}

func (h *Heap) Fix(i int) {
	heap.Fix(h, i)
}

func NewHeap(urls ...*TimedURL) *Heap {
	h := Heap(urls)
	heap.Init(&h)
	return &h
}
