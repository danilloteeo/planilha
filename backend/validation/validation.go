package validation

import (
	"fmt"
	"strings"

	"planly/backend/formula"
	"planly/backend/sheet"
	"planly/backend/templates"
)

const maxIssues = 500

type Issue struct {
	Row     int    `json:"row"`
	Col     int    `json:"col"`
	Cell    string `json:"cell"`
	Message string `json:"message"`
}

type Result struct {
	Issues []Issue `json:"issues"`
	Total  int     `json:"total"`
}

func (r *Result) OK() bool {
	return r.Total == 0
}

func (r *Result) add(row, col int, format string, args ...any) {
	r.Total++
	if len(r.Issues) < maxIssues {
		r.Issues = append(r.Issues, Issue{
			Row:     row,
			Col:     col,
			Cell:    sheet.CellName(row, col),
			Message: fmt.Sprintf(format, args...),
		})
	}
}

func Validate(s *sheet.Sheet, t *templates.Template) Result {
	res := Result{Issues: []Issue{}}
	for r := range s.RowCount() {
		for c, v := range s.Cells[r] {
			if v != "" {
				checkValue(&res, t, r, c, v)
			}
		}
	}
	if t.Layout == templates.Matrix {
		checkMatrix(&res, s, t)
	} else {
		checkList(&res, s, t)
	}
	return res
}

func checkValue(res *Result, t *templates.Template, r, c int, v string) {
	if formula.IsFormula(v) {
		if err := formula.Check(v); err != nil {
			res.add(r, c, "Fórmula inválida: %s.", err)
		}
		return
	}
	role := t.Role(r, c)
	switch role.Kind() {
	case templates.KindInteger:
		if _, ok := sheet.ParseInteger(v); !ok {
			res.add(r, c, "%s deve ser um número inteiro (valor: %q).", role.Label(), v)
		}
	case templates.KindDecimal, templates.KindCurrency:
		if _, ok := sheet.ParseNumber(v); !ok {
			res.add(r, c, "%s inválido(a): %q não é um número.", role.Label(), v)
		}
	}
}

func checkMatrix(res *Result, s *sheet.Sheet, t *templates.Template) {
	for _, f := range t.Fixed {
		if f.Role == templates.RoleCode && strings.TrimSpace(s.Get(f.Row, f.Col)) == "" {
			res.add(f.Row, f.Col, "Código da planilha obrigatório.")
		}
	}
	idRow := roleRow(t, templates.RoleClientID)
	nameRow := roleRow(t, templates.RoleClientName)
	clientIDs := map[string]int{}
	for c := t.FirstDataCol; c < s.ColCount(); c++ {
		if s.ColIsEmpty(c, 0) {
			continue
		}
		id := strings.TrimSpace(s.Get(idRow, c))
		if id == "" {
			res.add(idRow, c, "ID do cliente obrigatório na coluna %s.", sheet.ColName(c))
		} else if prev, dup := clientIDs[id]; dup {
			res.add(idRow, c, "ID do cliente %s repetido (também na coluna %s).", id, sheet.ColName(prev))
		} else {
			clientIDs[id] = c
		}
		if nameRow >= 0 && strings.TrimSpace(s.Get(nameRow, c)) == "" {
			res.add(nameRow, c, "Nome do cliente obrigatório na coluna %s.", sheet.ColName(c))
		}
	}
	checkProducts(res, s, t, []int{t.TotalCol})
}

func checkList(res *Result, s *sheet.Sheet, t *templates.Template) {
	for r := t.FirstDataRow; r < s.RowCount(); r++ {
		if s.RowIsEmpty(r) {
			continue
		}
		for c, col := range t.Columns {
			required := col.Role == templates.RoleSupplierID || col.Role == templates.RoleProductID || col.Role == templates.RoleQuantity
			if required && strings.TrimSpace(s.Get(r, c)) == "" {
				res.add(r, c, "Linha %d: campo %q obrigatório.", r+1, col.Header)
			}
		}
	}
}

func checkProducts(res *Result, s *sheet.Sheet, t *templates.Template, skip []int) {
	idCol, descCol := -1, -1
	for c, col := range t.Columns {
		switch col.Role {
		case templates.RoleProductID:
			idCol = c
		case templates.RoleProductDesc:
			descCol = c
		}
	}
	seen := map[string]int{}
	for r := t.FirstDataRow; r < s.RowCount(); r++ {
		if s.RowIsEmpty(r, skip...) {
			continue
		}
		id := strings.TrimSpace(s.Get(r, idCol))
		if id == "" {
			res.add(r, idCol, "Linha %d: ID do produto obrigatório.", r+1)
		} else if prev, dup := seen[id]; dup {
			res.add(r, idCol, "Linha %d: ID do produto %s repetido (também na linha %d).", r+1, id, prev+1)
		} else {
			seen[id] = r
		}
		if descCol >= 0 && strings.TrimSpace(s.Get(r, descCol)) == "" {
			res.add(r, descCol, "Linha %d: descrição do produto obrigatória.", r+1)
		}
	}
}

func roleRow(t *templates.Template, role templates.Role) int {
	for i, r := range t.HeaderRowRoles {
		if r == role {
			return i
		}
	}
	return -1
}
