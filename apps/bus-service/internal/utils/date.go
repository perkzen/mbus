package utils

import "time"

func Today() string {
	return time.Now().Format("2006-01-02")
}

func ValidateDate(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}

func nextOrTodayMatchingDate(match func(time.Weekday) bool) string {
	now := time.Now()
	for {
		if match(now.Weekday()) {
			return now.Format("2006-01-02")
		}
		now = now.AddDate(0, 0, 1)
	}
}

func Weekday() string {
	return nextOrTodayMatchingDate(func(d time.Weekday) bool {
		return d >= time.Monday && d <= time.Friday
	})
}

func Saturday() string {
	return nextOrTodayMatchingDate(func(d time.Weekday) bool {
		return d == time.Saturday
	})
}

func Sunday() string {
	return nextOrTodayMatchingDate(func(d time.Weekday) bool {
		return d == time.Sunday
	})
}
