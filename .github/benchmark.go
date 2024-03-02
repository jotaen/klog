package main

import (
	"fmt"
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/app"
	"github.com/jotaen/klog/klog/app/cli/terminalformat"
	"github.com/jotaen/klog/klog/parser"
	"math/rand"
	"os"
	"strconv"
	"time"
)

var serialiser = app.NewSerialiser(terminalformat.NewStyler(terminalformat.NO_COLOUR), false)

func main() {
	// Setup
	iterations, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	rand.Seed(int64(time.Now().Nanosecond()))

	// Generate records
	date := klog.Ɀ_Date_(0, 1, 1)
	for i := 0; i < iterations; i++ {
		date = date.PlusDays(1)
		r := klog.NewRecord(date)

		// Should total
		if i%2 == ri(0, 2) {
			r.SetShouldTotal(klog.NewDuration(ri(0, 23), ri(0, 59)))
		}

		// Summary
		text := rt(0, 5)
		if len(text) > 0 {
			r.SetSummary(klog.Ɀ_RecordSummary_(text...))
		}

		// Entries
		entriesCount := ri(1, 5)
		for j := 0; j < entriesCount; j++ {
			added := re()(&r)
			if !added {
				entriesCount++
			}
		}
		fmt.Println(parser.SerialiseRecords(serialiser, r).ToString())
	}
}

// ri = random integer
func ri(min int, max int) int {
	return rand.Intn(max+1-min) + min
}

// rt = random texts
func rt(rowsMin int, rowsMax int) []string {
	alphabet := "abcdefghijklmnopqrstuvwxyz"

	texts := make([]string, ri(rowsMin, rowsMax))
	for j := 0; j < len(texts); j++ {
		bs := make([]byte, ri(1, 50))
		for i := range bs {
			bs[i] = alphabet[ri(0, len(alphabet)-1)]
		}
		texts[j] = string(bs)
	}
	return texts
}

// re = random entry
func re() func(r *klog.Record) bool {
	text := rt(0, 2)
	var entrySummary klog.EntrySummary
	if len(text) > 0 {
		entrySummary = klog.Ɀ_EntrySummary_(text...)
	}
	entryAdders := []func(r *klog.Record) bool{
		func(r *klog.Record) bool {
			(*r).AddDuration(klog.NewDuration(ri(-2, 23), ri(0, 60)), entrySummary)
			return true
		},
		func(r *klog.Record) bool {
			(*r).AddRange(klog.Ɀ_Range_(
				klog.Ɀ_Time_(ri(0, 11), ri(0, 59)),
				klog.Ɀ_Time_(ri(12, 23), ri(0, 59)),
			), entrySummary)
			return true
		},
		func(r *klog.Record) bool {
			err := (*r).Start(klog.NewOpenRange(klog.Ɀ_Time_(ri(0, 23), ri(0, 59))), entrySummary)
			return err == nil
		},
	}
	return entryAdders[ri(0, len(entryAdders)-1)]
}
