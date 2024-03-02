package txt

// Error contains infos about a parsing error in a Line.
type Error interface {
	// Error is an alias for Message.
	Error() string

	// LineNumber returns the logical line number, as shown in an editor.
	LineNumber() int

	// LineText is the original text of the line.
	LineText() string

	// Position is the cursor position in the line, excluding the indentation.
	Position() int

	// Column is the cursor position in the line, including the indentation.
	Column() int

	// Length returns the number of erroneous characters.
	Length() int

	// Code returns a unique identifier of the error kind.
	Code() string

	// Title returns a short error description.
	Title() string

	// Details returns additional information, such as hints or further explanations.
	Details() string

	// Message is a combination of Title and Details.
	Message() string

	// Origin returns the origin of the error, such as the file name.
	Origin() string
	SetOrigin(string) Error
}

type err struct {
	context  Block
	origin   string
	line     int
	position int
	length   int
	code     string
	title    string
	details  string
}

func (e *err) Error() string                 { return e.Message() }
func (e *err) LineNumber() int               { return e.context.OverallLineIndex(e.line) + 1 }
func (e *err) LineText() string              { return e.context.Lines()[e.line].Text }
func (e *err) Position() int                 { return e.position }
func (e *err) Column() int                   { return e.position + 1 }
func (e *err) Length() int                   { return e.length }
func (e *err) Code() string                  { return e.code }
func (e *err) Title() string                 { return e.title }
func (e *err) Details() string               { return e.details }
func (e *err) Message() string               { return e.title + ": " + e.details }
func (e *err) Origin() string                { return e.origin }
func (e *err) SetOrigin(origin string) Error { e.origin = origin; return e }

func NewError(b Block, line int, start int, length int, code string, title string, details string) Error {
	return &err{
		context:  b,
		line:     line,
		position: start,
		length:   length,
		code:     code,
		title:    title,
		details:  details,
	}
}
