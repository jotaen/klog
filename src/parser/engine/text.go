package engine

type Text struct {
	Value            []rune
	PointerPosition  int
	LineNumber       int
	IndentationLevel int
}

var END_OF_TEXT int32 = -1

func (t *Text) Peek() rune {
	char := SubRune(t.Value, t.PointerPosition, 1)
	if char == nil {
		return END_OF_TEXT
	}
	return char[0]
}

func (t *Text) PeekUntil(isMatch func(rune) bool) Text {
	result := Text{
		PointerPosition: t.PointerPosition,
		Value:           nil,
		LineNumber:      t.LineNumber,
	}
	for i := t.PointerPosition; true; i++ {
		next := SubRune(t.Value, i, 1)
		if next == nil { // end of text
			return result
		}
		if isMatch(next[0]) {
			return result
		}
		result.Value = append(result.Value, next[0])
	}
	return result
}

func (t *Text) Advance(increment int) {
	t.PointerPosition += increment
}

func (t *Text) SkipWhitespace() {
	for t.Peek() == ' ' {
		t.Advance(1)
	}
	return
}

func (t *Text) Length() int {
	return len(t.Value)
}

func (t *Text) ToString() string {
	return string(t.Value)
}
