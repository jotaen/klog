package engine

import (
	"sync"
)

type dataIn[In any] struct {
	data In
	i    int
}

type dataOut[Out any, Err error] struct {
	data Out
	errs []Err
	i    int
}

type ParallelParser[In any, Out any, Err error] struct {
	NumberOfWorkers int
}

func (pp ParallelParser[In, Out, Err]) ParseAll(ins []In, parseOne func(In) (Out, []Err)) ([]Out, []Err) {
	outs := make([]Out, len(ins))
	var allErrs []Err

	// Set up channels.
	inChannel := make(chan dataIn[In])
	outChannel := make(chan dataOut[Out, Err])

	// Dispatch work.
	wg := &sync.WaitGroup{}
	wg.Add(len(ins))
	go func() {
		for i, in := range ins {
			inChannel <- dataIn[In]{in, i}
		}
		close(inChannel)
	}()

	// Spawn workers.
	for i := 0; i < pp.NumberOfWorkers; i++ {
		go func() {
			for in := range inChannel {
				data, errs := parseOne(in.data)
				out := dataOut[Out, Err]{i: in.i}
				if len(errs) > 0 {
					out.errs = errs
				} else {
					out.data = data
				}
				outChannel <- out
				wg.Done()
			}
		}()
	}

	go func() {
		wg.Wait()
		close(outChannel)
	}()

	for out := range outChannel {
		if out.errs == nil {
			outs[out.i] = out.data
		} else {
			allErrs = append(allErrs, out.errs...)
		}
	}

	if len(allErrs) > 0 {
		return nil, allErrs
	}
	return outs, nil
}
