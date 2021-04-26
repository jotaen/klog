package app

type Error interface {
	Error() string
	Details() string
	Original() error
}

type AppError struct {
	message  string
	details  string
	original error
}

func NewError(message string, details string, original error) Error {
	return AppError{message, details, original}
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
