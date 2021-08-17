package service

import . "klog"

type DayHash uint32

func NewDayHash(d Date) DayHash {
	hash := DayHash(0)                  // bit layout: ...YYYYYYYYYMMMMDDDDD
	hash = hash | DayHash(d.Day())<<0   // needs 5 bits max
	hash = hash | DayHash(d.Month())<<5 // needs 4 bits max
	hash = hash | DayHash(d.Year())<<9
	return hash
}
