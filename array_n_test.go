package lockfree_ds

import (
	"bytes"
	"sync"
	"testing"
)

var anpool = sync.Pool{New: func() any { return &ArrayN[*bytes.Buffer]{} }}

// call me using -race option
func TestArrayN(t *testing.T) {
	t.Run("", func(t *testing.T) {
		const N uint = 10

		var an ArrayN[*bytes.Buffer]
		an.Reserve(N)
		var wg sync.WaitGroup
		for i := 0; i < int(N); i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				buf, idx := an.Get()
				if buf == nil {
					buf = &bytes.Buffer{}
					an.Put(idx, buf)
				} else {
					buf.Reset()
				}
				buf.WriteString("foobar")
			}()
		}
		wg.Wait()
	})
}

// call me using -race option
func BenchmarkArrayN(b *testing.B) {
	b.Run("lock-free", func(b *testing.B) {
		b.ReportAllocs()
		const N uint = 1000
		for i := 0; i < b.N; i++ {
			an := anpool.Get().(*ArrayN[*bytes.Buffer])
			an.Reserve(N)
			var wg sync.WaitGroup // <--- memprofile: alloc
			for j := 0; j < int(N); j++ {
				wg.Add(1)
				go func() { // <--- memprofile: alloc
					defer wg.Done()
					buf, idx := an.Get()
					if buf == nil {
						buf = &bytes.Buffer{}
						an.Put(idx, buf)
					} else {
						buf.Reset()
					}
					buf.WriteString("foobar")
				}()
			}
			wg.Wait()
			an.Reset()
			anpool.Put(an)
		}
	})
	b.Run("mutex competitor", func(b *testing.B) {
		var anmpool = sync.Pool{New: func() any { return &anmux[*bytes.Buffer]{} }}

		b.ReportAllocs()
		const N uint = 1000
		for i := 0; i < b.N; i++ {
			anm := anmpool.Get().(*anmux[*bytes.Buffer])
			anm.reserve(N)
			var wg sync.WaitGroup // <--- memprofile: alloc
			for j := 0; j < int(N); j++ {
				wg.Add(1)
				go func() { // <--- memprofile: alloc
					defer wg.Done()
					buf, idx := anm.get()
					if buf == nil {
						buf = &bytes.Buffer{}
						anm.put(idx, buf)
					} else {
						buf.Reset()
					}
					buf.WriteString("foobar")
				}()
			}
			wg.Wait()
			anm.reset()
			anmpool.Put(anm)
		}
	})
}

// ArrayN mutex competitor.
type anmux[T any] struct {
	mux sync.Mutex
	idx uint32
	buf []T
}

func (a *anmux[T]) reserve(n uint) {
	a.idx = 0
	if a.buf == nil || cap(a.buf) < int(n) {
		a.buf = make([]T, n)
	}
	a.buf = a.buf[:n]
}

func (a *anmux[T]) get() (T, uint32) {
	a.mux.Lock()
	defer a.mux.Unlock()
	i := a.idx
	v := a.buf[a.idx]
	a.idx++
	return v, i
}

func (a *anmux[T]) put(idx uint32, value T) {
	a.mux.Lock()
	defer a.mux.Unlock()
	if idx < uint32(len(a.buf)) {
		a.buf[idx] = value
	}
}

func (a *anmux[T]) reset() {
	a.idx = 0
	a.buf = a.buf[:0]
}
