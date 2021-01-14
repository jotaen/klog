package record

func Total(r Record) Duration {
	total := NewDuration(0, 0)
	for _, e := range r.Entries() {
		switch v := e.Value().(type) {
		case Duration:
			total = total.Add(v)
			break
		case Range:
			total = total.Add(v.Duration())
			break
		}
	}
	return total
}

func Find(date Date, rs []Record) Record {
	for _, r := range rs {
		if r.Date() == date {
			return r
		}
	}
	return nil
}
