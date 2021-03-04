package json

import (
	"bytes"
	"encoding/json"
	. "klog"
	"klog/parser/parsing"
	"klog/service"
	"strings"
)

type Envelop struct {
	Records interface{} `json:"records"`
	Errors  interface{} `json:"errors"`
}

type EntryView struct {
	Type      string   `json:"type"`
	Summary   string   `json:"summary"`
	Tags      []string `json:"tags"`
	Total     string   `json:"total"`
	TotalMins int      `json:"total_mins"`
}

type OpenRangeView struct {
	EntryView
	Start     string `json:"start"`
	StartMins int    `json:"start_mins"`
}

type RangeView struct {
	OpenRangeView
	End     string `json:"end"`
	EndMins int    `json:"end_mins"`
}

type RecordView struct {
	Date            string        `json:"date"`
	Summary         string        `json:"summary"`
	Total           string        `json:"total"`
	TotalMins       int           `json:"total_mins"`
	ShouldTotal     string        `json:"should_total"`
	ShouldTotalMins int           `json:"should_total_mins"`
	Tags            []string      `json:"tags"`
	Entries         []interface{} `json:"entries"`
}

func ToJson(rs []Record, errs parsing.Errors) string {
	envelop := Envelop{
		Records: toView(rs),
		Errors:  errs,
	}
	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	enc.SetEscapeHTML(false)
	err := enc.Encode(&envelop)
	if err != nil {
		panic(err) // This should never happen
	}
	return strings.TrimRight(buffer.String(), "\n")
}

func toView(rs []Record) []RecordView {
	result := []RecordView{}
	for _, r := range rs {
		total := service.Total(r)
		v := RecordView{
			Date:            r.Date().ToString(),
			Summary:         r.Summary().ToString(),
			Total:           total.ToString(),
			TotalMins:       total.InMinutes(),
			ShouldTotal:     r.ShouldTotal().ToString(),
			ShouldTotalMins: r.ShouldTotal().InMinutes(),
			Tags:            r.Summary().Tags().ToStrings(),
			Entries:         toEntryViews(r.Entries()),
		}
		result = append(result, v)
	}
	return result
}

func toEntryViews(es []Entry) []interface{} {
	var views []interface{}
	for _, e := range es {
		base := EntryView{
			Summary:   e.Summary().ToString(),
			Tags:      e.Summary().Tags().ToStrings(),
			Total:     e.Duration().ToString(),
			TotalMins: e.Duration().InMinutes(),
		}
		if base.Tags == nil {
			base.Tags = []string{}
		}
		view := e.Unbox(func(r Range) interface{} {
			base.Type = "range"
			return RangeView{
				OpenRangeView: OpenRangeView{
					EntryView: base,
					Start:     r.Start().ToString(),
					StartMins: r.Start().MidnightOffset().InMinutes(),
				},
				End:     r.End().ToString(),
				EndMins: r.End().MidnightOffset().InMinutes(),
			}
		}, func(d Duration) interface{} {
			base.Type = "duration"
			return base
		}, func(o OpenRange) interface{} {
			base.Type = "open_range"
			return OpenRangeView{
				EntryView: base,
				Start:     o.Start().ToString(),
				StartMins: o.Start().MidnightOffset().InMinutes(),
			}
		})
		views = append(views, view)
	}
	return views
}
