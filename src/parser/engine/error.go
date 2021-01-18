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

type Error interface {
	Error() string
	Context() Text
	Position() int
	Length() int
	Code() string
	Message() string
	Set(string, string) Error
}

type err struct {
	context  Text
	position int
	length   int
	code     string
	message  string
}

func (e err) Error() string   { return e.message }
func (e err) Context() Text   { return e.context }
func (e err) Position() int   { return e.position }
func (e err) Length() int     { return e.length }
func (e err) Code() string    { return e.code }
func (e err) Message() string { return e.message }
func (e err) Set(code string, message string) Error {
	e.code = code
	e.message = message
	return e
}

func NewError(t Text, start int, length int) Error {
	return err{
		context:  t,
		position: start,
		length:   length,
	}
}
