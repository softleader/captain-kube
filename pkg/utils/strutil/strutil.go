package strutil

func Contains(vs []string, s string) bool {
	for _, v := range vs {
		if v == s {
			return true
		}
	}
	return false
}
