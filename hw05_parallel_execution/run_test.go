package hw05_parallel_execution //nolint:golint,stylecheck

import (
	"fmt"
	"math/rand"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
	"go.uber.org/goleak"
)

func TestRun(t *testing.T) {
	defer goleak.VerifyNone(t)

	t.Run("if were errors in first M tasks, than finished not more N+M tasks", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 10
		maxErrorsCount := 23
		result := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, ErrErrorsLimitExceeded, result)
		require.LessOrEqual(t, runTasksCount, int32(workersCount+maxErrorsCount), "extra tasks were started")
	})

	t.Run("tasks without errors", func(t *testing.T) {
		tasksCount := 50
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32
		var sumTime time.Duration

		for i := 0; i < tasksCount; i++ {
			taskSleep := time.Millisecond * time.Duration(rand.Intn(100))
			sumTime += taskSleep

			tasks = append(tasks, func() error {
				time.Sleep(taskSleep)
				atomic.AddInt32(&runTasksCount, 1)
				return nil
			})
		}

		workersCount := 5
		maxErrorsCount := 1

		start := time.Now()
		result := Run(tasks, workersCount, maxErrorsCount)
		elapsedTime := time.Since(start)
		require.Nil(t, result)

		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
		require.LessOrEqual(t, int64(elapsedTime), int64(sumTime/2), "tasks were run sequentially?")
	})

	t.Run("workers more than tasks", func(t *testing.T) {
		tasksCount := 5
		tasks := make([]Task, 0, tasksCount)

		for i := 0; i < tasksCount; i++ {
			tasks = append(tasks, func() error { return nil })
		}

		workersCount := 10
		maxErrorsCount := 0
		result := Run(tasks, workersCount, maxErrorsCount)

		require.Nil(t, result)
	})

	t.Run("ignore any errors", func(t *testing.T) {
		tasksCount := 100
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 3
		maxErrorsCount := -1
		result := Run(tasks, workersCount, maxErrorsCount)

		require.Nil(t, result)
		require.Equal(t, runTasksCount, int32(tasksCount), "not all tasks were completed")
	})

	t.Run("at least one error", func(t *testing.T) {
		tasksCount := 10
		tasks := make([]Task, 0, tasksCount)

		var runTasksCount int32

		for i := 0; i < tasksCount; i++ {
			err := fmt.Errorf("error from task %d", i)
			tasks = append(tasks, func() error {
				atomic.AddInt32(&runTasksCount, 1)
				return err
			})
		}

		workersCount := 1
		maxErrorsCount := 0
		result := Run(tasks, workersCount, maxErrorsCount)

		require.Equal(t, ErrErrorsLimitExceeded, result)
		require.LessOrEqual(t, runTasksCount, int32(1), "too many tasks completed")
	})

	t.Run("no tasks", func(t *testing.T) {
		tasks := make([]Task, 0, 0)
		workersCount := 1
		maxErrorsCount := 0
		result := Run(tasks, workersCount, maxErrorsCount)

		require.Nil(t, result)
	})

	t.Run("no workers", func(t *testing.T) {
		tasks := make([]Task, 0, 0)
		tasks = append(tasks, func() error { return nil })
		workersCount := 0
		maxErrorsCount := 0
		result := Run(tasks, workersCount, maxErrorsCount)

		require.Nil(t, result)
	})
}
