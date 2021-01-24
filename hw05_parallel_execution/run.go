package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"errors"
	"sync"
	"sync/atomic"
)

var ErrErrorsLimitExceeded = errors.New("errors limit exceeded")

type Task func() error

func Run(tasks []Task, n int, m int) error {
	var errCounter int32

	progress := make(chan Task, len(tasks))
	for _, t := range tasks {
		progress <- t
	}
	close(progress)

	wg := sync.WaitGroup{}
	for w := 0; w < n; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for {
				if m >= 0 && atomic.LoadInt32(&errCounter) > int32(m) {
					return
				}

				task, ok := <-progress

				if !ok {
					return
				}

				if err := task(); err != nil {
					atomic.AddInt32(&errCounter, 1)
				}
			}
		}()
	}
	wg.Wait()

	// all fine or ignore errors
	if errCounter == 0 || m < 0 {
		return nil
	}

	if m >= 0 && errCounter > int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
