package sheet

import (
	"strconv"
	"strings"
)

func ColName(c int) string {
	var b []byte
	for c >= 0 {
		b = append([]byte{byte('A' + c%26)}, b...)
		c = c/26 - 1
	}
	return string(b)
}

func ColIndex(name string) int {
	name = strings.ToUpper(strings.TrimSpace(name))
	if name == "" {
		return -1
	}
	n := 0
	for _, ch := range name {
		if ch < 'A' || ch > 'Z' {
			return -1
		}
		n = n*26 + int(ch-'A'+1)
	}
	return n - 1
}

func CellName(r, c int) string {
	return ColName(c) + strconv.Itoa(r+1)
}
