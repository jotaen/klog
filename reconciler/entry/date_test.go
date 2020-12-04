package entry

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOkayWhenValidDateIsPresent(t *testing.T) {
	date := Date{Year: 2020, Month: time.January, Day: 1}
	entry, err := Create(date)
	assert.Nil(t, err)
	assert.Equal(t, entry.Date(), date)
}

func TestErrorWhenDateIsInvalid(t *testing.T) {
	entry, err := Create(Date{Year: 2020, Month: time.January, Day: 99})
	assert.Nil(t, entry)
	assert.Equal(t, err.(*EntryError).Code, INVALID_DATE)
}
