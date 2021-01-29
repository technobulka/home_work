package hw06_pipeline_execution //nolint:golint,stylecheck
import (
	"sync"
)

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	out := make(Bi)

	go func() {
		wg := sync.WaitGroup{}

		var chains = make([]Bi, len(stages))
		for i := 0; i < len(stages); i++ {
			chains[i] = make(Bi)
		}

		for i, do := range stages {
			var next *Bi
			if i == len(stages)-1 {
				next = &out
			} else {
				next = &chains[i+1]
			}

			wg.Add(1)
			go func(i int, do Stage) {
				defer wg.Done()
				defer close(*next)

				for result := range do(chains[i]) {
					*next <- result
				}
			}(i, do)
		}

		for item := range in {
			select {
			case <-done:
				return
			case chains[0] <- item:
			}
		}
		close(chains[0])

		wg.Wait()
	}()

	return out
}
