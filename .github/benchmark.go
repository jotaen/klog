package main

import (
	"fmt"
	"github.com/jotaen/klog/klog"
	"github.com/jotaen/klog/klog/parser"
	"math/rand"
	"os"
	"strconv"
	"time"
)

func main() {
	// Setup
	iterations, err := strconv.Atoi(os.Args[1])
	if err != nil {
		panic(err)
	}
	rand.Seed(int64(time.Now().Nanosecond()))
	now := klog.NewDateFromGo(time.Now())

	// Generate records
	for i := 0; i < iterations; i++ {
		date := now.PlusDays((i + 1) * -1)
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
		fmt.Println(parser.SerialiseRecords(parser.PlainSerialiser{}, r).ToString())
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
			(*r).AddDuration(klog.NewDuration(ri(0, 23), ri(0, 60)), entrySummary)
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
			err := (*r).StartOpenRange(klog.Ɀ_Time_(ri(0, 23), ri(0, 59)), entrySummary)
			return err == nil
		},
	}
	return entryAdders[ri(0, len(entryAdders)-1)]
}
