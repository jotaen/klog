package app

import (
	"fmt"
	"github.com/jotaen/klog/klog/parser/txt"
)

type Code int

const (
	// GENERAL_ERROR should be used for generic, otherwise unspecified errors.
	GENERAL_ERROR Code = iota + 1

	// NO_INPUT_ERROR should be used if no input was specified.
	NO_INPUT_ERROR

	// NO_TARGET_FILE should be used if no target file was specified.
	NO_TARGET_FILE

	// IO_ERROR should be used for errors during I/O processes.
	IO_ERROR

	// CONFIG_ERROR should be used for config-folder-related problems.
	CONFIG_ERROR

	// NO_SUCH_BOOKMARK_ERROR should be used if the specified an unknown bookmark name.
	NO_SUCH_BOOKMARK_ERROR

	// NO_SUCH_FILE should be used if the specified file does not exit.
	NO_SUCH_FILE

	// LOGICAL_ERROR should be used syntax or logical violations.
	LOGICAL_ERROR
)

// ToInt returns the numeric value of the error. This is typically used as exit code.
func (c Code) ToInt() int {
	return int(c)
}

// Error is a representation of an application error.
type Error interface {
	// Error returns the error message.
	Error() string

	Is(error) bool

	// Details returns additional details, such as a hint how to solve the problem.
	Details() string

	// Original returns the original underlying error, if it exists.
	Original() error

	// Code returns the error code.
	Code() Code
}

type AppError struct {
	code     Code
	message  string
	details  string
	original error
}

func NewError(message string, details string, original error) Error {
	return NewErrorWithCode(GENERAL_ERROR, message, details, original)
}

func NewErrorWithCode(code Code, message string, details string, original error) Error {
	return AppError{code, message, details, original}
}

func (e AppError) Error() string {
	return e.message
}

func (e AppError) Is(err error) bool {
	_, ok := err.(AppError)
	return ok
}

func (e AppError) Details() string {
	return e.details
}

func (e AppError) Original() error {
	return e.original
}

func (e AppError) Code() Code {
	return e.code
}

type ParserErrors interface {
	Error
	All() []txt.Error
}

type parserErrors struct {
	errors []txt.Error
}

func NewParserErrors(errs []txt.Error) ParserErrors {
	return parserErrors{errs}
}

func (pe parserErrors) Error() string {
	return fmt.Sprintf("%d parsing error(s)", len(pe.errors))
}

func (e parserErrors) Is(err error) bool {
	_, ok := err.(parserErrors)
	return ok
}

func (pe parserErrors) Details() string {
	return fmt.Sprintf("%d parsing error(s)", len(pe.errors))
}

func (pe parserErrors) Original() error {
	return nil
}

func (pe parserErrors) Code() Code {
	return LOGICAL_ERROR
}

func (pe parserErrors) All() []txt.Error {
	return pe.errors
}
