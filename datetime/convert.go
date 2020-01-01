package datetime

import (
	"time"
)

// RoundDate ...
func RoundDate(d time.Time) time.Time {
	return time.Date(d.Year(), d.Month(), d.Day(), 0, 0, 0, 0, d.Location())
}

// RoundPartHour round by num a part of hour
// roundHour is 1 part, round30min is 2 part, round10min is 6 part
// return itself if 60 divisibility by num
func RoundPartHour(d time.Time, num int) time.Time {
	if 60%num != 0 {
		return d
	}

	m := 60 / num
	min := d.Minute() / m * m
	return time.Date(d.Year(), d.Month(), d.Day(), d.Hour(), min, 0, 0, d.Location())
}
