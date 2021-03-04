package parsing

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

type Error interface {
	Error() string
	Context() Line
	Position() int // text position _without_ indentation
	Column() int   // column number _with_ indentation
	Length() int
	Code() string
	Title() string
	Details() string
	Message() string
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
func (e *err) Column() int     { return len(e.context.originalIndentation) + e.position + 1 }
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

func NewError(t Line, start int, length int) Error {
	return &err{
		context:  t,
		position: start,
		length:   length,
	}
}
