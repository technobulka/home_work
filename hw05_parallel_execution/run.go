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

	ignoreErrors := m < 0
	progress := make(chan Task, len(tasks))

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

	for _, t := range tasks {
		progress <- t
	}
	close(progress)

	wg.Wait()

	if !ignoreErrors && errCounter > 0 && errCounter >= int32(m) {
		return ErrErrorsLimitExceeded
	}

	return nil
}
