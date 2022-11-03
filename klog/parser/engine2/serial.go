package engine2

type SerialParser[In any, Out any, Err error] struct{}

func (p SerialParser[In, Out, Err]) ParseAll(ins []In, parseOne func(In) (Out, []Err)) ([]Out, []Err) {
	outs := make([]Out, len(ins))
	var errs []Err
	for i, in := range ins {
		out, err := parseOne(in)
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
