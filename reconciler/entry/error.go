package entry

const (
	INVALID_DATE = "INVALID_DATE"
	NEGATIVE_TIME = "NEGATIVE_TIME"
)

type EntryError struct {
	Code string
}

func (e *EntryError) Error() string {
	return e.Code
}
