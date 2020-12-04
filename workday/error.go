package workday

const (
	INVALID_DATE  = "INVALID_DATE"
	NEGATIVE_TIME = "NEGATIVE_TIME"
)

type WorkDayError struct {
	Code string
}

func (e *WorkDayError) Error() string {
	return e.Code
}
