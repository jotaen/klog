package entry

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOkayWhenAddingValidTimes(t *testing.T) {
	date := Date{Year: 2020, Month: time.January, Day: 1}
	entry, _ := Create(date)
	err := entry.AddTime(Minutes(1))
	assert.Nil(t, err)
	assert.Equal(t, entry.Times(), []Minutes{Minutes(1)})
}

func TestErrorWhenAddingInvalidTimes(t *testing.T) {
	date := Date{Year: 2020, Month: time.January, Day: 1}
	entry, _ := Create(date)
	err := entry.AddTime(Minutes(-1))
	assert.Equal(t, err.(*EntryError).Code, NEGATIVE_TIME)
	assert.Equal(t, len(entry.Times()), 0)
}
