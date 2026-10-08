package sheet

import "slices"

type Sheet struct {
	Template  string     `json:"template"`
	Name      string     `json:"name"`
	SheetName string     `json:"sheetName"`
	Cells     [][]string `json:"cells"`
	ColWidths []float64  `json:"colWidths"`
}

func (s *Sheet) RowCount() int {
	return len(s.Cells)
}

func (s *Sheet) ColCount() int {
	n := 0
	for _, row := range s.Cells {
		if len(row) > n {
			n = len(row)
		}
	}
	if len(s.ColWidths) > n {
		n = len(s.ColWidths)
	}
	return n
}

func (s *Sheet) Get(r, c int) string {
	if r < 0 || r >= len(s.Cells) || c < 0 || c >= len(s.Cells[r]) {
		return ""
	}
	return s.Cells[r][c]
}

func (s *Sheet) Set(r, c int, v string) {
	if r < 0 || c < 0 {
		return
	}
	for len(s.Cells) <= r {
		s.Cells = append(s.Cells, nil)
	}
	row := s.Cells[r]
	if c >= len(row) {
		if v == "" {
			return
		}
		grown := make([]string, c+1)
		copy(grown, row)
		row = grown
	}
	row[c] = v
	s.Cells[r] = row
}

func (s *Sheet) Resize(rows, cols int) {
	for len(s.Cells) < rows {
		s.Cells = append(s.Cells, nil)
	}
	for r := range s.Cells {
		if len(s.Cells[r]) < cols {
			grown := make([]string, cols)
			copy(grown, s.Cells[r])
			s.Cells[r] = grown
		}
	}
	for len(s.ColWidths) < cols {
		s.ColWidths = append(s.ColWidths, 0)
	}
}

func (s *Sheet) RowIsEmpty(r int, skip ...int) bool {
	if r < 0 || r >= len(s.Cells) {
		return true
	}
	for c, v := range s.Cells[r] {
		if v == "" || slices.Contains(skip, c) {
			continue
		}
		return false
	}
	return true
}

func (s *Sheet) ColIsEmpty(c int, fromRow int) bool {
	for r := fromRow; r < len(s.Cells); r++ {
		if s.Get(r, c) != "" {
			return false
		}
	}
	return true
}
