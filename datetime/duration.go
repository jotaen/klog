package datetime

type Duration int // in minutes

func (d Duration) ToString() string {
	t := time{hour: int(int(d) / 60), minute: int(d) % 60}
	return t.ToString()
}

func CreateDurationFromString(hhmm string) (Duration, error) {
	t, err := CreateTimeFromString(hhmm)
	if err != nil {
		return 0, err
	}
	return Duration(t.Hour() * 60 + t.Minute()), nil
}
