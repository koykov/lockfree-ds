package lockfree_ds

import (
	"sync"
	"sync/atomic"
	"unsafe"
)

type MichaelScottQueue struct {
	head, tail unsafe.Pointer
	once       sync.Once
}

func (q *MichaelScottQueue) Enqueue(value any) {
	q.once.Do(q.init)

	e := &entryQ{payload: value}
	for {
		tail := atomic.LoadPointer(&q.tail)
		next := atomic.LoadPointer(&(*entryQ)(tail).next)

		// check tail has been updated in parallel thread
		if tail == atomic.LoadPointer(&q.tail) {
			if next == nil {
				// update next of tail
				if atomic.CompareAndSwapPointer(&(*entryQ)(tail).next, next, unsafe.Pointer(e)) {
					// try to update tail
					atomic.CompareAndSwapPointer(&q.tail, tail, unsafe.Pointer(e))
					return
				}
			} else {
				// let parallel thread update tail
				atomic.CompareAndSwapPointer(&q.tail, tail, unsafe.Pointer(e))
			}
		}
	}
}

func (q *MichaelScottQueue) Dequeue() (any, bool) {
	q.once.Do(q.init)

	for {
		head := atomic.LoadPointer(&q.head)
		tail := atomic.LoadPointer(&q.tail)
		hnext := atomic.LoadPointer(&(*entryQ)(head).next)

		// check head didn't update in parallel thread
		if head != atomic.LoadPointer(&q.head) {
			continue
		}

		if head == tail {
			// queue may be empty (head and tail point to empty entry)
			if hnext == nil {
				// queue is empty
				return nil, false
			} else {
				atomic.CompareAndSwapPointer(&q.tail, tail, hnext)
			}
		} else {
			// queue isn't empty, try dequeue
			value := (*entryQ)(hnext).payload
			if atomic.CompareAndSwapPointer(&q.head, head, hnext) {
				return value, true
			}
		}
	}
}

func (q *MichaelScottQueue) init() {
	empty := unsafe.Pointer(&entryQ{})
	q.head = empty
	q.tail = empty
}

type entryQ struct {
	payload any
	next    unsafe.Pointer
}
