package app

const (
	GENERAL_ERROR = iota + 1
	NO_INPUT_ERROR
	NO_TARGET_FILE
	IO_ERROR
	BOOKMARK_ERROR
)

type Error interface {
	Error() string
	Details() string
	Original() error
	Code() int
}

type AppError struct {
	code     int
	message  string
	details  string
	original error
}

func NewError(message string, details string, original error) Error {
	return NewErrorWithCode(GENERAL_ERROR, message, details, original)
}

func NewErrorWithCode(code int, message string, details string, original error) Error {
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

func (e AppError) Code() int {
	return e.code
}
