package engine

import "fmt"

type Errors interface {
	Get() []Error
	Error() string
}

type errors struct {
	errors []Error
}

func NewErrors(errs []Error) Errors {
	return errors{errs}
}

func (pe errors) Error() string {
	return fmt.Sprintf("%d parsing errors", len(pe.errors))
}

func (pe errors) Get() []Error {
	return pe.errors
}

// Error contains infos about a parsing error in a Line.
type Error interface {
	// Error is an alias for Message.
	Error() string

	// Context is the Line in which the error occurred.
	Context() Line

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

	// Set fills in Code, Title and Message (in this order).
	Set(string, string, string) Error
}

type err struct {
	context  Line
	position int
	length   int
	code     string
	title    string
	details  string
}

func (e *err) Error() string   { return e.Message() }
func (e *err) Context() Line   { return e.context }
func (e *err) Position() int   { return e.position }
func (e *err) Column() int     { return e.position + 1 }
func (e *err) Length() int     { return e.length }
func (e *err) Code() string    { return e.code }
func (e *err) Title() string   { return e.title }
func (e *err) Details() string { return e.details }
func (e *err) Message() string { return e.title + ": " + e.details }
func (e *err) Set(code string, title string, details string) Error {
	e.code = code
	e.title = title
	e.details = details
	return e
}

func NewError(t Line, start int, length int, code string, title string, details string) Error {
	return &err{
		context:  t,
		position: start,
		length:   length,
		code:     code,
		title:    title,
		details:  details,
	}
}
