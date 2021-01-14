package record

import "klog/datetime"

func TotalWorkTime(r Record) datetime.Duration {
	total := datetime.NewDuration(0, 0)
	for _, t := range r.Durations() {
		total = total.Add(t)
	}
	for _, r := range r.Ranges() {
		total = total.Add(r.Duration())
	}
	return total
}
