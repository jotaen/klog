package app

type Error interface {
	Error() string
	Help() string
}

type appError struct {
	message string
	help    string
}

func (e appError) Error() string {
	return e.message
}

func (e appError) Help() string {
	return e.help
}
