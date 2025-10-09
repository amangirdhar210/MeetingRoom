package utils

import "time"

func IsTimeRangeValid(start, end time.Time) bool {
	return start.Before(end)
}

func Overlaps(start1, end1, start2, end2 time.Time) bool {
	return start1.Before(end2) && start2.Before(end1)
}
