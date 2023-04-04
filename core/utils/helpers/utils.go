package helpers

import (
	"os"
	"strconv"
)

// GetEnv return env value for the given environment variable, if it does exists, return fallback
func GetEnv(key, fallback string) string {
	value := os.Getenv(key)
	if len(value) == 0 {
		return fallback
	}
	return value
}

// ValidPageNumbers will evaluate if the given page number or item_per_page is valid
func ValidPageNumbers(pgVal string, name string) (bool, uint64) {
	PageValue, ok := strconv.ParseUint(pgVal, 10, 16)
	if ok != nil || PageValue < 1 {
		return true, 0
	}
	return false, PageValue
}
