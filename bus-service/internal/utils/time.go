package utils

import (
	"fmt"
	"time"
)

func ParseClock(s string) (time.Time, error) {
	return time.Parse("15:04", s)
}

func FormatDuration(from, to time.Time) string {
	dur := to.Sub(from)
	return fmt.Sprintf("%02d:%02d", int(dur.Hours()), int(dur.Minutes())%60)
}
