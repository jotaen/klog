package engine

import (
	"sync"
)

type ParallelBatchParser[Txt any, Int any, Out any, Err any] struct {
	SerialParser    SerialParser[Txt, Int, Out, Err]
	NumberOfWorkers int
}

type batchResult[Out any, Err any] struct {
	outs []Out
	errs []Err
}

func (p ParallelBatchParser[Txt, Int, Out, Err]) Parse(txt Txt) ([]Out, []Err) {
	allInts := p.SerialParser.PreProcess(txt)
	// Batch up and dispatch to workers.
	wg := &sync.WaitGroup{}
	batches := splitUp(allInts, p.NumberOfWorkers)
	wg.Add(len(batches))
	resultChannel := make(chan batchResult[Out, Err])
	for _, b := range batches {
		go func(ints []Int) {
			defer wg.Done()
			outs, errs := p.SerialParser.parseAll(ints)
			resultChannel <- batchResult[Out, Err]{outs, errs}
		}(b)
	}

	// Wait for workers to finish.
	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	// Collect results.
	i := 0
	outs := make([]Out, len(allInts))
	var allErrs []Err
	for result := range resultChannel {
		if len(result.errs) > 0 {
			allErrs = append(allErrs, result.errs...)
		} else {
			for _, o := range result.outs {
				outs[i] = o
				i++
			}
		}
	}
	if len(allErrs) > 0 {
		return nil, allErrs
	}
	return outs, nil
}

func splitUp[T any](slice []T, batchSize int) [][]T {
	batches := make([][]T, 0, (len(slice)+batchSize-1)/batchSize)
	for batchSize < len(slice) {
		slice, batches = slice[batchSize:], append(batches, slice[0:batchSize:batchSize])
	}
	return append(batches, slice)
}
