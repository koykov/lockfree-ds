package lockfree_ds

import (
	"context"
	"math/rand"
	"sync"
	"testing"
	"time"
)

// call me using -race option
func TestTreiberStack(t *testing.T) {
	t.Run("", func(t *testing.T) {
		var wg sync.WaitGroup
		var stk TreiberStack
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		for i := 0; i < 10; i++ {
			wg.Add(1)
			go func(ctx context.Context, stk *TreiberStack) {
				defer wg.Done()
				for {
					select {
					case <-ctx.Done():
						return
					default:
						v := rand.Intn(100)
						if v < 50 {
							stk.Push(v)
						} else {
							if x, ok := stk.Pop(); ok {
								if x.(int) >= 50 {
									t.Fail()
									return
								}
							}
						}
					}
				}
			}(ctx, &stk)
		}
		_ = cancel
		<-ctx.Done()
	})
}
