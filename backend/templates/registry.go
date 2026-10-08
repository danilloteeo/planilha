package templates

import (
	"strings"
	"unicode/utf8"
)

var registry = []*Template{preco(), venda(), compra()}

func All() []*Template {
	return registry
}

func Get(id string) (*Template, bool) {
	for _, t := range registry {
		if t.ID == id {
			return t, true
		}
	}
	return nil, false
}

func Detect(cells [][]string) (*Template, bool) {
	for _, t := range registry {
		if t.Matches(cells) {
			return t, true
		}
	}
	return nil, false
}

func SanitizeSheetName(name, fallback string) string {
	name = strings.TrimSpace(strings.Map(func(r rune) rune {
		if strings.ContainsRune(`[]:*?/\`, r) || r < 32 {
			return -1
		}
		return r
	}, name))
	name = strings.Trim(name, "'")
	for utf8.RuneCountInString(name) > 31 {
		_, size := utf8.DecodeLastRuneInString(name)
		name = name[:len(name)-size]
	}
	if name == "" {
		return fallback
	}
	return name
}
