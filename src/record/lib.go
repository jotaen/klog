package record

import "klog/datetime"

func Total(r Record) datetime.Duration {
	total := datetime.NewDuration(0, 0)
	for _, e := range r.Entries() {
		total = total.Add(e.Total())
	}
	return total
}

func Find(date datetime.Date, rs []Record) Record {
	for _, r := range rs {
		if r.Date() == date {
			return r
		}
	}
	return nil
}
