package env

import (
	"os"
	"strconv"
)

func Lookup(key, defaultValue string) string {
	if v, found := os.LookupEnv(key); found {
		return v
	}
	return defaultValue
}

func LookupInt(key string, defaultValue int) int {
	if v, found := os.LookupEnv(key); found {
		if i, err := strconv.Atoi(v); err != nil {
			return i
		}
	}
	return defaultValue
}
