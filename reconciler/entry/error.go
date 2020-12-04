package entry

const (
	INVALID_DATE = "The date does not represent a valid day in the calendar"
	NEGATIVE_TIME = "A time cannot be a negative value"
)

type EntryError struct {
	Code string
}

func (e *EntryError) Error() string {
	return e.Code
}
