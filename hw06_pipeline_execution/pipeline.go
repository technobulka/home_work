package hw06_pipeline_execution //nolint:golint,stylecheck

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	for _, stage := range stages {
		out := make(Bi)

		go func(in In, out Bi) {
			defer close(out)

			for {
				select {
				case <-done:
					return
				case val, ok := <-in:
					if !ok {
						return
					}
					out <- val
				}
			}
		}(in, out)

		in = stage(out)
	}

	return in
}
