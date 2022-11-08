package engine

import "github.com/jotaen/klog/klog/parser/txt"

type SerialParser[T any] struct {
	ParseOne func(txt.Block) (T, []txt.Error)
}

func (p SerialParser[T]) Parse(text string) ([]T, []txt.Block, []txt.Error) {
	ts, blocks, _, errs, hasErrors := p.mapParse(text)
	if hasErrors {
		return nil, nil, flatten[txt.Error](errs)
	}
	return ts, blocks, nil
}

// mapParse parses the text. All 3 return arrays have the same arity, and the last
// bool indicates whether any errors occurred.
func (p SerialParser[T]) mapParse(text string) ([]T, []txt.Block, int, [][]txt.Error, bool) {
	var ts []T
	var blocks []txt.Block
	var errs [][]txt.Error
	totalBytesConsumed := 0
	totalLines := 0
	hasErrors := false
	for {
		block, bytesConsumed := txt.ParseBlock(text[totalBytesConsumed:], totalLines)
		if bytesConsumed == 0 || block == nil {
			break
		}
		totalLines += len(block.Lines())
		totalBytesConsumed += bytesConsumed
		t, err := p.ParseOne(block)
		ts = append(ts, t)
		blocks = append(blocks, block)
		errs = append(errs, err)
		if err != nil {
			hasErrors = true
		}
	}
	return ts, blocks, totalBytesConsumed, errs, hasErrors
}

func flatten[T any](xss [][]T) []T {
	var result []T
	for _, xs := range xss {
		if len(xs) == 0 {
			continue
		}
		result = append(result, xs...)
	}
	return result
}
