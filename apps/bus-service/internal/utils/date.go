package utils

import "time"

func Today() string {
	return time.Now().Format("2006-01-02")
}

func ValidateDate(dateStr string) bool {
	_, err := time.Parse("2006-01-02", dateStr)
	return err == nil
}
