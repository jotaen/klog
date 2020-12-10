package parser

type ParserError struct {
	code     string
	location string
}

type errorCollection struct {
	collection []ParserError
}

func (p *errorCollection) add(err ParserError) {
	p.collection = append(p.collection, err)
}

func parserError(code string, location string) ParserError {
	return ParserError{
		code:     code,
		location: location,
	}
}

func fromError(err error, location string) ParserError {
	return ParserError{
		code:     err.Error(),
		location: location,
	}
}

func (e ParserError) Error() string {
	return "" // TODO compose good error message
}

func ToErrors(parserErrors []ParserError) []error {
	errs := []error{}
	for _, e := range parserErrors {
		errs = append(errs, e)
	}
	return errs
}
