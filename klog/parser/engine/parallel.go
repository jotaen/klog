package engine

import (
	"github.com/jotaen/klog/klog/parser/txt"
	"math"
	"sync"
	"unicode/utf8"
)

type ParallelBatchParser[T any] struct {
	SerialParser[T]
	NumberOfWorkers int
}

type batchResult[T any] struct {
	index    int
	headText string
	values   []T
	tailText string
	blocks   []txt.Block
	errs     [][]txt.Error
}

func (p ParallelBatchParser[T]) Parse(text string) ([]T, []txt.Block, []txt.Error) {
	if p.NumberOfWorkers <= 0 {
		panic("ILLEGAL_WORKER_SIZE")
	}
	batches := splitIntoChunks(text, p.NumberOfWorkers)
	allResults := p.processAsync(batches, func(batchIndex int, batchText string) batchResult[T] {
		result := batchResult[T]{batchIndex, "", nil, "", nil, nil}

		if len(batchText) == 0 {
			return result
		}

		_, headBytesConsumed := txt.ParseBlock(batchText, 1)
		result.headText = batchText[:headBytesConsumed]
		if len(batchText) == headBytesConsumed { // The entire batchText was a single block
			return result
		}

		batchText = batchText[headBytesConsumed:]
		values, blocks, bytesConsumed, errs, _ := p.SerialParser.parse(batchText)
		if len(blocks) == 0 { // The remainder was empty or all blank
			result.tailText = batchText
		} else { // The remainder was more than one block
			result.values = values[:len(values)-1]
			result.blocks = blocks[:len(blocks)-1]
			result.errs = errs[:len(errs)-1]
			result.tailText = batchText[bytesConsumed-countBytes(blocks[len(blocks)-1]):]
		}

		return result
	})

	// Process remainders and flatten results.
	var allValues []T
	var allBlocks []txt.Block
	var allErrs []txt.Error
	carryText := ""
	for _, result := range allResults {
		carryText += result.headText
		if len(result.blocks) > 0 {
			carryValues, carryBlocks, _, carryErrs, hasErrors := p.SerialParser.parse(carryText)
			allValues = append(allValues, carryValues...)
			allBlocks = append(allBlocks, carryBlocks...)
			if hasErrors {
				allErrs = append(allErrs, flatten(carryErrs)...)
			}
			carryText = ""
			allValues = append(allValues, result.values...)
			allBlocks = append(allBlocks, result.blocks...)
			allErrs = append(allErrs, flatten(result.errs)...)
		}
		carryText += result.tailText
	}
	carryValues, carryBlocks, _, carryErrs, hasErrors := p.SerialParser.parse(carryText)
	allValues = append(allValues, carryValues...)
	allBlocks = append(allBlocks, carryBlocks...)
	lineCount := 0
	for _, b := range allBlocks {
		b.SetPrecedingLineCount(lineCount)
		lineCount += len(b.Lines())
	}
	if hasErrors {
		allErrs = append(allErrs, flatten(carryErrs)...)
	}
	if len(allErrs) > 0 {
		return nil, nil, allErrs
	}
	return allValues, allBlocks, nil
}

func (p ParallelBatchParser[T]) processAsync(batches []string, work func(int, string) batchResult[T]) []batchResult[T] {
	wg := &sync.WaitGroup{}
	wg.Add(len(batches))
	resultChannel := make(chan batchResult[T])
	for i, b := range batches {
		go func(batchIndex int, batchText string) {
			defer wg.Done()
			result := work(batchIndex, batchText)
			resultChannel <- result
		}(i, b)
	}

	// Wait for workers to finish.
	go func() {
		wg.Wait()
		close(resultChannel)
	}()

	// Collect results.
	allResults := make([]batchResult[T], len(batches))
	for result := range resultChannel {
		allResults[result.index] = result
	}

	return allResults
}

// splitIntoChunks divides a string into n substrings of roughly equal byte-size
// (not character-count). The chunk’s byte size might differ slightly: (a) because
// the last chunk contains the remainder, which will probably be smaller, and (b)
// because the chunks are never divided in between UTF-8 code points.
func splitIntoChunks(txt string, numberOfBatches int) []string {
	batchByteSize := int(math.Ceil(float64(len(txt)) / float64(numberOfBatches)))
	batches := make([]string, numberOfBatches)
	pointer := 0
	for i := 0; i < numberOfBatches; i++ {
		nextPointer := pointer + batchByteSize
		for nextPointer < len(txt) && !utf8.RuneStart(txt[nextPointer]) {
			nextPointer++
		}
		if nextPointer > len(txt) {
			batches[i] = txt[pointer:]
			break
		} else {
			batches[i] = txt[pointer:nextPointer]
		}
		pointer = nextPointer
	}
	return batches
}

func countBytes(b txt.Block) int {
	result := 0
	for _, l := range b.Lines() {
		result += len(l.Original())
	}
	return result
}
