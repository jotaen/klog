package engine

type SerialParser[Txt any, Int any, Out any, Err any] struct {
	PreProcess func(Txt) []Int
	ParseOne   func(Int) (Out, []Err)
}

func (p SerialParser[Txt, Int, Out, Err]) Parse(text Txt) ([]Out, []Err) {
	return p.parseAll(p.PreProcess(text))
}

func (p SerialParser[Txt, Int, Out, Err]) parseAll(ints []Int) ([]Out, []Err) {
	outs := make([]Out, len(ints))
	var errs []Err
	for i, in := range ints {
		out, err := p.ParseOne(in)
		if err != nil {
			errs = append(errs, err...)
			continue
		}
		outs[i] = out
	}
	if errs != nil {
		return nil, errs
	}
	return outs, errs
}
