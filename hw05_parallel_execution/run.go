package hw05parallelexecution

import (
	"errors"
	"sync"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

// Run starts tasks in n goroutines and stops its work when receiving m errors from tasks.
func Run(tasks []Task, n, m int) error {
	if m <= 0 {
		return ErrErrorsLimitExceeded
	}

	var errCount int
	var wg sync.WaitGroup
	var mx sync.Mutex

	ch := make(chan Task, len(tasks))

	for _, task := range tasks {
		ch <- task
	}
	close(ch)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for {
				mx.Lock()
				if errCount >= m {
					mx.Unlock()
					return
				}
				mx.Unlock()

				task, ok := <-ch
				if !ok {
					return
				}

				if err := task(); err == nil {
					continue
				}
				mx.Lock()
				errCount++
				mx.Unlock()
			}
		}()
	}
	wg.Wait()
	if errCount >= m {
		return ErrErrorsLimitExceeded
	}
	return nil
}
