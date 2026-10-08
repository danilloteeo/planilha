package formula

import (
	"errors"
	"fmt"
	"strings"
	"unicode"
)

func IsFormula(v string) bool {
	return len(v) > 1 && v[0] == '='
}

type dialect struct {
	sepIn, sepOut         rune
	decimalIn, decimalOut rune
	lookup                func(string) (string, bool)
}

var (
	toExcel   = dialect{';', ',', ',', '.', LookupPT}
	fromExcel = dialect{',', ';', '.', ',', LookupEN}
)

func ToExcel(v string) string {
	return translate(strings.TrimPrefix(v, "="), toExcel)
}

func FromExcel(v string) string {
	return "=" + translate(strings.TrimPrefix(v, "="), fromExcel)
}

func translate(expr string, d dialect) string {
	src := []rune(expr)
	var out strings.Builder
	for i := 0; i < len(src); {
		ch := src[i]
		switch {
		case ch == '"':
			j, _ := skipString(src, i)
			out.WriteString(string(src[i:j]))
			i = j
		case ch == '\'':
			j := skipQuoted(src, i, '\'')
			out.WriteString(string(src[i:j]))
			i = j
		case unicode.IsDigit(ch):
			j := i
			for j < len(src) && unicode.IsDigit(src[j]) {
				j++
			}
			out.WriteString(string(src[i:j]))
			if j+1 < len(src) && src[j] == d.decimalIn && unicode.IsDigit(src[j+1]) {
				out.WriteRune(d.decimalOut)
				j++
				k := j
				for k < len(src) && unicode.IsDigit(src[k]) {
					k++
				}
				out.WriteString(string(src[j:k]))
				j = k
			}
			i = j
		case isWordStart(ch):
			j := i
			for j < len(src) && isWordPart(src[j]) {
				j++
			}
			word := string(src[i:j])
			if next := nextNonSpace(src, j); next < len(src) && src[next] == '(' {
				if mapped, ok := d.lookup(word); ok {
					word = mapped
				}
			}
			out.WriteString(word)
			i = j
		case ch == d.sepIn:
			out.WriteRune(d.sepOut)
			i++
		default:
			out.WriteRune(ch)
			i++
		}
	}
	return out.String()
}

func Check(v string) error {
	if !IsFormula(v) {
		return errors.New("a fórmula deve começar com =")
	}
	src := []rune(strings.TrimPrefix(v, "="))
	if strings.TrimSpace(string(src)) == "" {
		return errors.New("fórmula vazia")
	}
	depth := 0
	for i := 0; i < len(src); {
		ch := src[i]
		switch {
		case ch == '"':
			j, closed := skipString(src, i)
			if !closed {
				return errors.New("texto sem aspas de fechamento")
			}
			i = j
			continue
		case ch == '(':
			depth++
		case ch == ')':
			depth--
			if depth < 0 {
				return errors.New("parêntese fechado sem abertura")
			}
		case isWordStart(ch):
			j := i
			for j < len(src) && isWordPart(src[j]) {
				j++
			}
			word := string(src[i:j])
			if next := nextNonSpace(src, j); next < len(src) && src[next] == '(' {
				if _, ok := LookupPT(word); !ok {
					if _, ok := LookupEN(word); !ok {
						return fmt.Errorf("função desconhecida: %s", word)
					}
				}
			}
			i = j
			continue
		}
		i++
	}
	if depth != 0 {
		return errors.New("parênteses não balanceados")
	}
	last := lastNonSpace(src)
	if last >= 0 && strings.ContainsRune("+-*/^&;,(=<>", src[last]) {
		return errors.New("fórmula incompleta")
	}
	return nil
}

func skipString(src []rune, i int) (int, bool) {
	j := i + 1
	for j < len(src) {
		if src[j] == '"' {
			if j+1 < len(src) && src[j+1] == '"' {
				j += 2
				continue
			}
			return j + 1, true
		}
		j++
	}
	return len(src), false
}

func skipQuoted(src []rune, i int, q rune) int {
	j := i + 1
	for j < len(src) && src[j] != q {
		j++
	}
	if j < len(src) {
		j++
	}
	return j
}

func isWordStart(ch rune) bool {
	return unicode.IsLetter(ch) || ch == '_' || ch == '$'
}

func isWordPart(ch rune) bool {
	return unicode.IsLetter(ch) || unicode.IsDigit(ch) || ch == '_' || ch == '.' || ch == '$'
}

func nextNonSpace(src []rune, i int) int {
	for i < len(src) && unicode.IsSpace(src[i]) {
		i++
	}
	return i
}

func lastNonSpace(src []rune) int {
	i := len(src) - 1
	for i >= 0 && unicode.IsSpace(src[i]) {
		i--
	}
	return i
}
