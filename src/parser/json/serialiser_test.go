package json

import (
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/engine"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSerialiseEmptyRecords(t *testing.T) {
	json := ToJson([]Record{}, nil, false)
	assert.Equal(t, `{"records":[],"errors":null}`, json)
}

func TestSerialiseEmptyArrayIfNoErrors(t *testing.T) {
	json := ToJson(nil, nil, false)
	assert.Equal(t, `{"records":[],"errors":null}`, json)
}

func TestSerialisePrettyPrinted(t *testing.T) {
	json := ToJson(nil, nil, true)
	assert.Equal(t, `{
  "records": [],
  "errors": null
}`, json)
}

func TestSerialiseMinimalRecord(t *testing.T) {
	json := ToJson(func() []Record {
		r := NewRecord(Ɀ_Date_(2000, 12, 31))
		return []Record{r}
	}(), nil, false)
	assert.Equal(t, `{"records":[{`+
		`"date":"2000-12-31",`+
		`"summary":"",`+
		`"total":"0m",`+
		`"total_mins":0,`+
		`"should_total":"0m",`+
		`"should_total_mins":0,`+
		`"diff":"0m",`+
		`"diff_mins":0,`+
		`"tags":[],`+
		`"entries":[]`+
		`}],"errors":null}`, json)
}

func TestSerialiseFullBlownRecord(t *testing.T) {
	json := ToJson(func() []Record {
		r := NewRecord(Ɀ_Date_(2000, 12, 31))
		r.SetSummary(Ɀ_RecordSummary_("Hello #World"))
		r.SetShouldTotal(NewDuration(7, 30))
		r.AddDuration(NewDuration(2, 3), NewEntrySummary("#some #thing"))
		r.AddRange(Ɀ_Range_(Ɀ_TimeYesterday_(23, 44), Ɀ_Time_(5, 23)), nil)
		r.StartOpenRange(Ɀ_TimeTomorrow_(0, 28), NewEntrySummary("Started #todo"))
		return []Record{r}
	}(), nil, false)
	assert.Equal(t, `{"records":[{`+
		`"date":"2000-12-31",`+
		`"summary":"Hello #World",`+
		`"total":"7h42m",`+
		`"total_mins":462,`+
		`"should_total":"7h30m!",`+
		`"should_total_mins":450,`+
		`"diff":"+12m",`+
		`"diff_mins":12,`+
		`"tags":["#world"],`+
		`"entries":[{`+
		`"type":"duration",`+
		`"summary":"#some #thing",`+
		`"tags":["#some","#thing"],`+
		`"total":"2h3m",`+
		`"total_mins":123`+
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

func TestSerialiseParserErrors(t *testing.T) {
	json := ToJson(nil, []engine.Error{
		parser.ErrorInvalidDate().New(engine.Line{
			Text:       "2018-99-99",
			LineNumber: 7,
		}, 0, 10),
	}, false)
	assert.Equal(t, `{"records":null,"errors":[{`+
		`"line":7,`+
		`"column":1,`+
		`"length":10,`+
		`"title":"Invalid date",`+
		`"details":"Please make sure that the date format is either YYYY-MM-DD or YYYY/MM/DD, and that its value represents a valid day in the calendar."`+
		`}]}`, json)
}
