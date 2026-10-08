package templates

import (
	"strconv"
	"strings"

	"planly/backend/sheet"
)

type Layout string

const (
	Matrix Layout = "matrix"
	List   Layout = "list"
)

type Role string

const (
	RoleEmpty        Role = "empty"
	RoleTitle        Role = "title"
	RoleCode         Role = "code"
	RoleHeader       Role = "header"
	RoleRoute        Role = "route"
	RoleClientID     Role = "clientId"
	RoleClientName   Role = "clientName"
	RoleProductID    Role = "productId"
	RoleProductDesc  Role = "productDesc"
	RoleSupplierID   Role = "supplierId"
	RoleSupplierName Role = "supplierName"
	RoleTotal        Role = "total"
	RolePrice        Role = "price"
	RoleQuantity     Role = "quantity"
)

type Kind string

const (
	KindText     Kind = "text"
	KindInteger  Kind = "integer"
	KindDecimal  Kind = "decimal"
	KindCurrency Kind = "currency"
	KindFormula  Kind = "formula"
)

func (r Role) Kind() Kind {
	switch r {
	case RoleCode, RoleRoute, RoleClientID, RoleProductID, RoleSupplierID:
		return KindInteger
	case RoleQuantity:
		return KindDecimal
	case RolePrice:
		return KindCurrency
	case RoleTotal:
		return KindFormula
	default:
		return KindText
	}
}

func (r Role) Label() string {
	switch r {
	case RoleTitle:
		return "Título"
	case RoleCode:
		return "Código"
	case RoleHeader:
		return "Cabeçalho"
	case RoleRoute:
		return "Rota"
	case RoleClientID:
		return "ID do cliente"
	case RoleClientName:
		return "Cliente"
	case RoleProductID:
		return "ID do produto"
	case RoleProductDesc:
		return "Produto"
	case RoleSupplierID:
		return "ID do fornecedor"
	case RoleSupplierName:
		return "Fornecedor"
	case RoleTotal:
		return "Soma"
	case RolePrice:
		return "Preço"
	case RoleQuantity:
		return "Quantidade"
	default:
		return ""
	}
}

type Style struct {
	Bold   bool    `json:"bold"`
	Size   float64 `json:"size"`
	HAlign string  `json:"hAlign"`
	VAlign string  `json:"vAlign"`
	Border bool    `json:"border"`
	Wrap   bool    `json:"wrap"`
	NumFmt string  `json:"numFmt"`
}

type Column struct {
	Header string  `json:"header"`
	Role   Role    `json:"role"`
	Width  float64 `json:"width"`
}

type FixedCell struct {
	Row      int    `json:"row"`
	Col      int    `json:"col"`
	Text     string `json:"text"`
	Role     Role   `json:"role"`
	ReadOnly bool   `json:"readOnly"`
}

type RowHeight struct {
	Row    int     `json:"row"`
	Height float64 `json:"height"`
}

type Template struct {
	ID               string           `json:"id"`
	Label            string           `json:"label"`
	Description      string           `json:"description"`
	Layout           Layout           `json:"layout"`
	DefaultSheetName string           `json:"defaultSheetName"`
	DefaultFileName  string           `json:"defaultFileName"`
	FirstDataRow     int              `json:"firstDataRow"`
	FirstDataCol     int              `json:"firstDataCol"`
	HeaderRowRoles   []Role           `json:"headerRowRoles"`
	Columns          []Column         `json:"columns"`
	DataRole         Role             `json:"dataRole"`
	DataWidth        float64          `json:"dataWidth"`
	TotalCol         int              `json:"totalCol"`
	Fixed            []FixedCell      `json:"fixed"`
	FreezeRow        int              `json:"freezeRow"`
	FreezeCol        int              `json:"freezeCol"`
	DefaultRowHeight float64          `json:"defaultRowHeight"`
	RowHeights       []RowHeight      `json:"rowHeights"`
	Styles           map[string]Style `json:"styles"`
	DefaultClients   int              `json:"defaultClients"`
	DefaultProducts  int              `json:"defaultProducts"`
}

type NewOptions struct {
	Name      string `json:"name"`
	SheetName string `json:"sheetName"`
	Clients   int    `json:"clients"`
	Products  int    `json:"products"`
	Code      string `json:"code"`
}

func (t *Template) fixedAt(r, c int) (FixedCell, bool) {
	for _, f := range t.Fixed {
		if f.Row == r && f.Col == c {
			return f, true
		}
	}
	return FixedCell{}, false
}

func (t *Template) Role(r, c int) Role {
	if f, ok := t.fixedAt(r, c); ok {
		return f.Role
	}
	if t.Layout == Matrix {
		if r < t.FirstDataRow {
			if c >= t.FirstDataCol && r < len(t.HeaderRowRoles) {
				return t.HeaderRowRoles[r]
			}
			return RoleEmpty
		}
		if c < t.FirstDataCol {
			if c < len(t.Columns) {
				return t.Columns[c].Role
			}
			return RoleEmpty
		}
		return t.DataRole
	}
	if r < t.FirstDataRow || c >= len(t.Columns) {
		return RoleEmpty
	}
	return t.Columns[c].Role
}

func (t *Template) ReadOnly(r, c int) bool {
	if f, ok := t.fixedAt(r, c); ok {
		return f.ReadOnly
	}
	return t.Role(r, c) == RoleTotal
}

func (t *Template) StyleFor(role Role) Style {
	if st, ok := t.Styles[string(role)]; ok {
		return st
	}
	return t.Styles["default"]
}

func (t *Template) ColWidth(c int) float64 {
	if c < len(t.Columns) {
		return t.Columns[c].Width
	}
	if t.Layout == Matrix {
		return t.DataWidth
	}
	return 9.14
}

func (t *Template) New(opts NewOptions) *sheet.Sheet {
	s := &sheet.Sheet{
		Template:  t.ID,
		Name:      strings.TrimSpace(opts.Name),
		SheetName: SanitizeSheetName(opts.SheetName, t.DefaultSheetName),
	}
	if s.Name == "" {
		s.Name = t.Label
	}
	rows, cols := t.FirstDataRow, len(t.Columns)
	if t.Layout == Matrix {
		rows += clamp(opts.Products, 0, 100000)
		cols = t.FirstDataCol + clamp(opts.Clients, 0, 5000)
	} else {
		rows += clamp(opts.Products, 0, 100000)
	}
	s.Resize(rows, cols)
	for c := range cols {
		s.ColWidths[c] = t.ColWidth(c)
	}
	for _, f := range t.Fixed {
		s.Set(f.Row, f.Col, f.Text)
	}
	if code := strings.TrimSpace(opts.Code); code != "" {
		for _, f := range t.Fixed {
			if f.Role == RoleCode {
				s.Set(f.Row, f.Col, code)
			}
		}
	}
	t.Normalize(s)
	return s
}

func (t *Template) Normalize(s *sheet.Sheet) {
	s.Template = t.ID
	minCols := len(t.Columns)
	if t.Layout == Matrix {
		minCols = t.FirstDataCol
	}
	rows := max(s.RowCount(), t.FirstDataRow)
	s.Resize(rows, max(s.ColCount(), minCols))
	for c := range s.ColWidths {
		if s.ColWidths[c] <= 0 {
			s.ColWidths[c] = t.ColWidth(c)
		}
	}
	for _, f := range t.Fixed {
		if f.ReadOnly {
			s.Set(f.Row, f.Col, f.Text)
		}
	}
	if t.TotalCol >= 0 {
		for r := t.FirstDataRow; r < s.RowCount(); r++ {
			s.Set(r, t.TotalCol, t.TotalFormula(r, s.ColCount()))
		}
	}
}

func (t *Template) TotalFormula(r, cols int) string {
	if cols <= t.FirstDataCol {
		return ""
	}
	row := r + 1
	return "=SOMA(" + sheet.ColName(t.FirstDataCol) + strconv.Itoa(row) + ":" + sheet.ColName(cols-1) + strconv.Itoa(row) + ")"
}

func (t *Template) Matches(cells [][]string) bool {
	get := func(r, c int) string {
		if r < len(cells) && c < len(cells[r]) {
			return strings.TrimSpace(cells[r][c])
		}
		return ""
	}
	matched := 0
	for _, f := range t.Fixed {
		if !f.ReadOnly {
			continue
		}
		if !strings.EqualFold(get(f.Row, f.Col), f.Text) {
			return false
		}
		matched++
	}
	return matched > 0
}

func clamp(v, lo, hi int) int {
	return min(max(v, lo), hi)
}
