package utils

import (
	"log"
	"time"
)

func Retry(description string, maxAttempts int, delay time.Duration, fn func() error) error {
	var err error

	for i := 1; i <= maxAttempts; i++ {
		err = fn()
		if err == nil {
			log.Printf("✅ %s succeeded on attempt %d", description, i)
			return nil
		}

		log.Printf("⚠️ %s failed on attempt %d/%d: %v", description, i, maxAttempts, err)
		time.Sleep(delay)
	}

	log.Printf("❌ %s failed after %d attempts. Last error: %v", description, maxAttempts, err)
	return err
}
