package lockfree_ds

import (
	"math"
	"sync/atomic"
)

type ArrayN[T any] struct {
	idx uint32
	buf []T
}

func (a *ArrayN[T]) Reserve(n uint) {
	a.idx = math.MaxUint32
	if a.buf == nil || cap(a.buf) < int(n) {
		a.buf = make([]T, n)
	}
	a.buf = a.buf[:n]
}

func (a *ArrayN[T]) Get() (T, uint32) {
	i := atomic.AddUint32(&a.idx, 1)
	return a.buf[i], i
}

func (a *ArrayN[T]) Put(idx uint32, value T) {
	if idx < uint32(len(a.buf)) {
		a.buf[idx] = value
	}
}

func (a *ArrayN[T]) Reset() {
	a.idx = math.MaxUint32
	a.buf = a.buf[:0]
}
