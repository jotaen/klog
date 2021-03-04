package json

import (
	"github.com/stretchr/testify/assert"
	. "klog"
	"testing"
)

func TestSerialiseEmptyRecords(t *testing.T) {
	json := ToJson([]Record{}, nil)
	assert.Equal(t, `{"records":[],"errors":null}`, json)
}

func TestSerialiseEmptyArrayIfNoErrors(t *testing.T) {
	json := ToJson(nil, nil)
	assert.Equal(t, `{"records":[],"errors":null}`, json)
}

func TestSerialiseFullBlownRecord(t *testing.T) {
	json := ToJson(func() []Record {
		r := NewRecord(Ɀ_Date_(2000, 12, 31))
		r.SetSummary("Hello #World")
		r.SetShouldTotal(NewDuration(7, 30))
		r.AddDuration(NewDuration(1, 30), "#some #thing")
		return []Record{r}
	}(), nil)
	assert.Equal(t, `{"records":[{`+
		`"date":"2000-12-31",`+
		`"summary":"Hello #World",`+
		`"total":"1h30m",`+
		`"total_mins":90,`+
		`"should_total":"7h30m!",`+
		`"should_total_mins":450,`+
		`"tags":["#world"],`+
		`"entries":[{`+
		`"type":"duration",`+
		`"summary":"#some #thing",`+
		`"tags":["#some","#thing"],`+
		`"total":"1h30m",`+
		`"total_mins":90`+
		`}]`+
		`}],"errors":null}`, json)
}
