/*
Package json contains the logic of serialising Record’s as JSON.
*/
package json

import (
	"bytes"
	"encoding/json"
	. "github.com/jotaen/klog/src"
	"github.com/jotaen/klog/src/parser"
	"github.com/jotaen/klog/src/parser/parsing"
	"github.com/jotaen/klog/src/service"
	"strings"
)

// ToJson serialises records into their JSON representation. The output
// structure is RecordView at the top level.
func ToJson(rs []Record, errs parsing.Errors, prettyPrint bool) string {
	envelop := func() Envelop {
		if errs == nil {
			return Envelop{
				Records: toRecordViews(rs),
				Errors:  nil,
			}
		} else {
			return Envelop{
				Records: nil,
				Errors:  toErrorViews(errs),
			}
		}
	}()
	buffer := new(bytes.Buffer)
	enc := json.NewEncoder(buffer)
	if prettyPrint {
		enc.SetIndent("", "  ")
	}
	enc.SetEscapeHTML(false)
	err := enc.Encode(&envelop)
	if err != nil {
		panic(err) // This should never happen
	}
	return strings.TrimRight(buffer.String(), "\n")
}

func toRecordViews(rs []Record) []RecordView {
	result := []RecordView{}
	for _, r := range rs {
		total := service.Total(r)
		should := r.ShouldTotal()
		diff := service.Diff(should, total)
		v := RecordView{
			Date:            r.Date().ToString(),
			Summary:         parser.SummaryText(r.Summary()).ToString(),
			Total:           total.ToString(),
			TotalMins:       total.InMinutes(),
			ShouldTotal:     should.ToString(),
			ShouldTotalMins: should.InMinutes(),
			Diff:            diff.ToStringWithSign(),
			DiffMins:        diff.InMinutes(),
			Tags:            toTagViews(r.Summary().Tags()),
			Entries:         toEntryViews(r.Entries()),
		}
		result = append(result, v)
	}
	return result
}

func toTagViews(ts TagSet) []string {
	result := ts.ToStrings()
	if result == nil {
		return []string{}
	}
	return result
}

func toEntryViews(es []Entry) []interface{} {
	views := []interface{}{}
	for _, e := range es {
		base := EntryView{
			Summary:   parser.SummaryText(e.Summary()).ToString(),
			Tags:      toTagViews(e.Summary().Tags()),
			Total:     e.Duration().ToString(),
			TotalMins: e.Duration().InMinutes(),
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

func toErrorViews(errs parsing.Errors) []ErrorView {
	var result []ErrorView
	for _, e := range errs.Get() {
		result = append(result, ErrorView{
			Line:    e.Context().LineNumber,
			Column:  e.Column(),
			Length:  e.Length(),
			Title:   e.Title(),
			Details: e.Details(),
		})
	}
	return result
}
