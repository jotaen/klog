package app

type Error interface {
	Error() string
	Details() string
}

type AppError struct {
	message string
	details string
}

func NewError(message string, details string) Error {
	return AppError{message, details}
}

func (e AppError) Error() string {
	return e.message
}

func (e AppError) Details() string {
	return e.details
}
