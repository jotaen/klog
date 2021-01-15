package engine

import "errors"

type Line struct {
	Pos  int // Current position while parsing
	Text []rune
	Nr   int
}

func (l *Line) Peek(length int) string {
	return substr(l.Text, l.Pos, length)
}

func (l *Line) PeekUntil(char rune) (string, error) {
	result := ""
	for i := l.Pos; true; i++ {
		next := substr(l.Text, i, 1)
		if next == string(char) {
			return result, nil
		}
		result += next
		if next == "\n" {
			break
		}
	}
	return result, errors.New("EXPECTED_CHARACTER_NOT_FOUND")
}

func (l *Line) Advance(increment int) {
	l.Pos += increment
}

func (l *Line) SkipWhitespace() {
	for l.Peek(1) == " " {
		l.Advance(1)
	}
	return
}
