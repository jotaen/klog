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
		r.AddRange(Ɀ_Range_(Ɀ_TimeYesterday_(23, 44), Ɀ_Time_(5, 23)), "")
		r.StartOpenRange(Ɀ_TimeTomorrow_(0, 28), "Started #todo")
		return []Record{r}
	}(), nil)
	assert.Equal(t, `{"records":[{`+
		`"date":"2000-12-31",`+
		`"summary":"Hello #World",`+
		`"total":"7h9m",`+
		`"total_mins":429,`+
		`"should_total":"7h30m!",`+
		`"should_total_mins":450,`+
		`"tags":["#world"],`+
		`"entries":[{`+
		`"type":"duration",`+
		`"summary":"#some #thing",`+
		`"tags":["#some","#thing"],`+
		`"total":"1h30m",`+
		`"total_mins":90`+
		`},{`+
		`"type":"range",`+
		`"summary":"",`+
		`"tags":[],`+
		`"total":"5h39m",`+
		`"total_mins":339,`+
		`"start":"<23:44",`+
		`"start_mins":-16,`+
		`"end":"5:23",`+
		`"end_mins":323`+
		`},{`+
		`"type":"open_range",`+
		`"summary":"Started #todo",`+
		`"tags":["#todo"],`+
		`"total":"0m",`+
		`"total_mins":0,`+
		`"start":"0:28>",`+
		`"start_mins":1468`+
		`}]`+
		`}],"errors":null}`, json)
}
