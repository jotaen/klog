package engine

import "github.com/jotaen/klog/klog/parser/txt"

type SerialParser[T any] struct {
	ParseOne func(txt.Block) (T, []txt.Error)
}

func (p SerialParser[T]) Parse(text string) ([]T, []txt.Block, []txt.Error) {
	blocks := txt.GroupIntoBlocks(text)
	result := make([]T, len(blocks))
	var errs []txt.Error
	for i, in := range blocks {
		out, err := p.ParseOne(in)
		if err != nil {
			errs = append(errs, err...)
			continue
		}
		result[i] = out
	}
	if errs != nil {
		return nil, nil, errs
	}
	return result, blocks, errs
}
