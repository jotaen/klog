package entry

import (
	"cloud.google.com/go/civil"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestOkayWhenValidDateIsPresent(t *testing.T) {
	day := Entry{
		Date: civil.Date{Year: 2020, Month: time.January, Day: 3},
	}
	assert.Nil(t, day.Check())
}

func TestErrorWhenDateIsNotInitialised(t *testing.T) {
	day := Entry{}
	assert.Contains(t, day.Check(), EntryError{ Code: INVALID_DATE })
}

func TestErrorWhenDateIsInvalid(t *testing.T) {
	day := Entry{
		Date: civil.Date{Year: 2020, Month: time.January, Day: 99},
	}
	assert.Contains(t, day.Check(), EntryError{ Code: INVALID_DATE })
}

func TestErrorWhenTimesAreInvalid(t *testing.T) {
	day := Entry{
		Times: []Minutes{Minutes(1), Minutes(-5)},
	}
	assert.Contains(t, day.Check(), EntryError{ Code: NEGATIVE_TIME })
}
