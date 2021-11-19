package app

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

	// CONFIG_ERROR should be used for .klog-folder-related problems.
	CONFIG_ERROR

	// NO_SUCH_BOOKMARK_ERROR should be used if the specified an unknown bookmark name.
	NO_SUCH_BOOKMARK_ERROR

	// NO_SUCH_FILE should be used if the specified file does not exit.
	NO_SUCH_FILE
)

// ToInt returns the numeric value of the error. This is typically used as exit code.
func (c Code) ToInt() int {
	return int(c)
}

// Error is a representation of an application error.
type Error interface {
	// Error returns the error message.
	Error() string

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

func (e AppError) Details() string {
	return e.details
}

func (e AppError) Original() error {
	return e.original
}

func (e AppError) Code() Code {
	return e.code
}
