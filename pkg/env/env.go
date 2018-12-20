package env

import (
	"os"
	"strconv"
)

// 從 os.LookupEnv 找系統參數, 如果沒找到則返回 defaultValue
func Lookup(key, defaultValue string) string {
	if v, found := os.LookupEnv(key); found {
		return v
	}
	return defaultValue
}

// 從 os.LookupEnv 找系統參數, 如果沒找到或找到後在轉換成 int 時發生錯誤, 則返回 defaultValue
func LookupInt(key string, defaultValue int) int {
	if v, found := os.LookupEnv(key); found {
		if i, err := strconv.Atoi(v); err == nil {
			return i
		}
	}
	return defaultValue
}

// 從 os.LookupEnv 找系統參數, 如果沒找到或找到後在轉換成 bool 時發生錯誤, 則返回 defaultValue
func LookupBool(key string, defaultValue bool) bool {
	if v, found := os.LookupEnv(key); found {
		if b, err := strconv.ParseBool(v); err == nil {
			return b
		}
	}
	return defaultValue
}
