package strutil

import (
	"github.com/pmezard/go-difflib/difflib"
	"strings"
)

// Contains 回傳 s 是否包含於 vs 中
func Contains(vs []string, s string) bool {
	for _, v := range vs {
		if v == s {
			return true
		}
	}
	return false
}

// DiffNewLines 回傳 a 跟 b 不同樣的行數
func DiffNewLines(a, b string) []string {
	diff := difflib.UnifiedDiff{
		A:       difflib.SplitLines(a),
		B:       difflib.SplitLines(b),
		Context: 3,
		Eol:     "\n",
	}

	diffResult, _ := difflib.GetUnifiedDiffString(diff)
	lines := strings.Split(diffResult, "\n")

	// 取出新增的
	var result []string
	for _, line := range lines {
		if len(line) > 1 && line[0] == '+' {
			result = append(result, line[1:])
		}
	}

	return result
}
