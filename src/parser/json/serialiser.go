package json

import (
	"encoding/json"
	. "klog"
	"klog/parser/parsing"
	"klog/service"
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

type RecordView struct {
	Date            string      `json:"date"`
	Summary         string      `json:"summary"`
	Total           string      `json:"total"`
	TotalMins       int         `json:"total_mins"`
	ShouldTotal     string      `json:"should_total"`
	ShouldTotalMins int         `json:"should_total_mins"`
	Tags            []string    `json:"tags"`
	Entries         []EntryView `json:"entries"`
}

func ToJson(rs []Record, errs parsing.Errors) string {
	envelop := Envelop{
		Records: toView(rs),
		Errors:  errs,
	}
	result, err := json.Marshal(&envelop)
	if err != nil {
		panic(err) // This should never happen
	}
	return string(result)
}

func toView(rs []Record) []RecordView {
	result := []RecordView{}
	for _, r := range rs {
		entries := func() []EntryView {
			var evs []EntryView
			for _, e := range r.Entries() {
				entryType := e.Unbox(func(r Range) interface{} {
					return "range"
				}, func(d Duration) interface{} {
					return "duration"
				}, func(openRange OpenRange) interface{} {
					return "open_range"
				}).(string)
				evs = append(evs, EntryView{
					Type:      entryType,
					Summary:   e.Summary().ToString(),
					Tags:      e.Summary().Tags().ToStrings(),
					Total:     e.Duration().ToString(),
					TotalMins: e.Duration().InMinutes(),
				})
				return evs
			}
			return evs
		}()
		total := service.Total(r)
		v := RecordView{
			Date:            r.Date().ToString(),
			Summary:         r.Summary().ToString(),
			Total:           total.ToString(),
			TotalMins:       total.InMinutes(),
			ShouldTotal:     r.ShouldTotal().ToString(),
			ShouldTotalMins: r.ShouldTotal().InMinutes(),
			Tags:            r.Summary().Tags().ToStrings(),
			Entries:         entries,
		}
		result = append(result, v)
	}
	return result
}
