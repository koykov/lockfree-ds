package lockfree_ds

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"
)

func TestMichaelScottQueue(t *testing.T) {
	t.Run("", func(t *testing.T) {
		var wg sync.WaitGroup
		var q MichaelScottQueue
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		for i := 0; i < 1; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					default:
						n := rand.Intn(100)
						if n < 50 {
							q.Enqueue(n)
						} else {
							if m, ok := q.Dequeue(); ok {
								if m.(int) >= 50 {
									t.Fail()
									return
								}
							}
						}
					}
				}
			}()
		}
		wg.Wait()
		_ = cancel
	})
}
