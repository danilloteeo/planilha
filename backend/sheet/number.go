package sheet

import (
	"math"
	"regexp"
	"strconv"
	"strings"
)

var thousandsOnly = regexp.MustCompile(`^-?\d{1,3}(\.\d{3})+$`)

func ParseNumber(s string) (float64, bool) {
	s = strings.TrimSpace(s)
	s = strings.TrimPrefix(s, "R$")
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, false
	}
	percent := strings.HasSuffix(s, "%")
	s = strings.TrimSuffix(s, "%")
	switch {
	case strings.Contains(s, ","):
		s = strings.ReplaceAll(s, ".", "")
		s = strings.Replace(s, ",", ".", 1)
	case thousandsOnly.MatchString(s):
		s = strings.ReplaceAll(s, ".", "")
	}
	f, err := strconv.ParseFloat(s, 64)
	if err != nil || math.IsNaN(f) || math.IsInf(f, 0) {
		return 0, false
	}
	if percent {
		f /= 100
	}
	return f, true
}

func ParseInteger(s string) (int64, bool) {
	f, ok := ParseNumber(s)
	if !ok || f != math.Trunc(f) || math.Abs(f) > 1e15 {
		return 0, false
	}
	return int64(f), true
}

func FormatNumber(f float64) string {
	return strings.Replace(strconv.FormatFloat(f, 'f', -1, 64), ".", ",", 1)
}

func FromRawNumber(raw string) (string, bool) {
	f, err := strconv.ParseFloat(strings.TrimSpace(raw), 64)
	if err != nil {
		return raw, false
	}
	return FormatNumber(roundTo(f, 10)), true
}

func roundTo(f float64, digits int) float64 {
	p := math.Pow(10, float64(digits))
	return math.Round(f*p) / p
}
