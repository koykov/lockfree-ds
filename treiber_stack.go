package lockfree_ds

import (
	"sync/atomic"
	"unsafe"
)

type TreiberStack struct {
	head unsafe.Pointer
}

func (s *TreiberStack) Push(value any) {
	e := entry{payload: value}
	for {
		headptr := atomic.LoadPointer(&s.head)
		e.next = headptr
		if atomic.CompareAndSwapPointer(&s.head, headptr, unsafe.Pointer(&e)) {
			break
		}
	}
}

func (s *TreiberStack) Pop() (value any, ok bool) {
	for {
		headptr := atomic.LoadPointer(&s.head)
		if headptr == nil {
			return
		}
		headentry := (*entry)(headptr)
		next := atomic.LoadPointer(&headentry.next)
		if ok = atomic.CompareAndSwapPointer(&s.head, headptr, next); ok {
			value = headentry.payload
			return
		}
	}
}

type entry struct {
	payload any
	next    unsafe.Pointer
}
