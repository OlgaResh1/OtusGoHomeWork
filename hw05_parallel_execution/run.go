package hw05parallelexecution

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 || n <= 0 {
		return ErrErrorsLimitExceeded
	}

	var errCount int32
	var wg sync.WaitGroup

	ch := make(chan Task)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for task, ok := <-ch; ok; task, ok = <-ch {
				if err := task(); err != nil {
					atomic.AddInt32(&errCount, 1)
				}
			}
		}()
	}
	for _, task := range tasks {
		ch <- task

		if atomic.LoadInt32(&errCount) >= int32(m) {
			break
		}
	}
	close(ch)

	wg.Wait()
	if atomic.LoadInt32(&errCount) >= int32(m) {
		return ErrErrorsLimitExceeded
	}
	return nil
}
