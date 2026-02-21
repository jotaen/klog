package json

import (
	"testing"

	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/parser"
	"github.com/jotaen/klog/klog/parser/txt"
	"github.com/stretchr/testify/assert"
)

func TestSerialiseEmptyRecords(t *testing.T) {
	json := ToJson([]klog.Record{}, nil, nil, false)
	assert.Equal(t, `{"records":[],"warnings":null,"errors":null}`, json)
}

func TestSerialiseEmptyArrayIfNoErrors(t *testing.T) {
	json := ToJson(nil, nil, nil, false)
	assert.Equal(t, `{"records":[],"warnings":null,"errors":null}`, json)
}

func TestSerialisePrettyPrinted(t *testing.T) {
	json := ToJson(nil, nil, nil, true)
	assert.Equal(t, `{
  "records": [],
  "warnings": null,
  "errors": null
}`, json)
}

func TestSerialiseMinimalRecord(t *testing.T) {
	json := ToJson(func() []klog.Record {
		r := klog.NewRecord(klog.Ɀ_Date_(2000, 12, 31))
		return []klog.Record{r}
	}(), nil, nil, false)
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
		`}],"warnings":null,"errors":null}`, json)
}

func TestSerialiseFullBlownRecord(t *testing.T) {
	json := ToJson(func() []klog.Record {
		r := klog.NewRecord(klog.Ɀ_Date_(2000, 12, 31))
		r.SetSummary(klog.Ɀ_RecordSummary_("Hello #World", "What’s up?"))
		r.SetShouldTotal(klog.NewDuration(7, 30))
		r.AddDuration(klog.NewDuration(2, 3), klog.Ɀ_EntrySummary_("#some #thing"))
		r.AddRange(klog.Ɀ_Range_(klog.Ɀ_TimeYesterday_(23, 44), klog.Ɀ_Time_(5, 23)), nil)
		r.Start(klog.NewOpenRange(klog.Ɀ_TimeTomorrow_(0, 28)), klog.Ɀ_EntrySummary_("Started #todo=nr4", "still on #it"))
		return []klog.Record{r}
	}(), nil, nil, false)
	assert.Equal(t, `{"records":[{`+
		`"date":"2000-12-31",`+
		`"summary":"Hello #World\nWhat’s up?",`+
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
		`"summary":"Started #todo=nr4\nstill on #it",`+
		`"tags":["#it","#todo=nr4"],`+
		`"total":"0m",`+
		`"total_mins":0,`+
		`"start":"0:28>",`+
		`"start_mins":1468`+
		`}]`+
		`}],"warnings":null,"errors":null}`, json)
}

func TestSerialiseParserErrors(t *testing.T) {
	block, _ := txt.ParseBlock("2018-99-99\n asdf", 6)
	json := ToJson(nil, []txt.Error{
		parser.ErrorInvalidDate().New(block, 0, 0, 10),
		parser.ErrorMalformedSummary().New(block, 1, 3, 5).SetOrigin("/a/b/c/file.klg"),
	}, nil, false)
	assert.Equal(t, `{"records":null,"warnings":null,"errors":[{`+
		`"line":7,`+
		`"column":1,`+
		`"length":10,`+
		`"title":"Invalid date",`+
		`"details":"Please make sure that the date format is either YYYY-MM-DD or YYYY/MM/DD, and that its value represents a valid day in the calendar.",`+
		`"file":""`+
		`},{`+
		`"line":8,`+
		`"column":4,`+
		`"length":5,`+
		`"title":"Malformed summary",`+
		`"details":"Summary lines cannot start with blank characters, such as non-breaking spaces.",`+
		`"file":"/a/b/c/file.klg"`+
		`}]}`, json)
}

func TestSerialiseWarnings(t *testing.T) {
	t.Run("include warnings if records are ok", func(t *testing.T) {
		json := ToJson(nil, nil, []string{"Caution!", "Beware!"}, false)
		assert.Equal(t, `{"records":[],"warnings":["Caution!","Beware!"],"errors":null}`, json)
	})
	t.Run("ignore warnings if records are not ok", func(t *testing.T) {
		block, _ := txt.ParseBlock("2018-99-99", 6)
		json := ToJson(nil, []txt.Error{
			parser.ErrorInvalidDate().New(block, 0, 0, 10),
		}, []string{"Caution!", "Beware!"}, false)
		assert.Equal(t, `{"records":null,"warnings":null,"errors":[{`+
			`"line":7,`+
			`"column":1,`+
			`"length":10,`+
			`"title":"Invalid date",`+
			`"details":"Please make sure that the date format is either YYYY-MM-DD or YYYY/MM/DD, and that its value represents a valid day in the calendar.",`+
			`"file":""`+
			`}]}`, json)
	})
}
