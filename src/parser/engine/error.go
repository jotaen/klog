package engine

type Error struct {
	Context  Text
	Position int
	Length   int
	Code     string
	Message  string
}

func (e Error) Error() string {
	return e.Message
}

func NewError(t Text, start int, length int) Error {
	return Error{
		Context:  t,
		Position: start,
		Length:   length,
	}
}
