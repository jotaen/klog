package app

type Code int

const (
	GENERAL_ERROR Code = iota + 1
	NO_INPUT_ERROR
	NO_TARGET_FILE
	IO_ERROR
	BOOKMARK_ERROR
)

func (c Code) ToInt() int {
	return int(c)
}

type Error interface {
	Error() string
	Details() string
	Original() error
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
