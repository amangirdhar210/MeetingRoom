package utils

func IsTimeRangeValid(start, end int64) bool {
	return start < end
}

func Overlaps(start1, end1, start2, end2 int64) bool {
	return start1 < end2 && start2 < end1
}
