package excel

import (
	"context"
	"fmt"
	"io"

	"github.com/xuri/excelize/v2"

	"planly/backend/formula"
	"planly/backend/sheet"
	"planly/backend/templates"
)

type Progress func(done, total int)

const progressStep = 2000

func ExportXLSX(ctx context.Context, s *sheet.Sheet, t *templates.Template, path string, progress Progress) error {
	f := excelize.NewFile()
	defer f.Close()

	name := templates.SanitizeSheetName(s.SheetName, t.DefaultSheetName)
	if err := f.SetSheetName("Sheet1", name); err != nil {
		return fmt.Errorf("renomear aba: %w", err)
	}
	sw, err := f.NewStreamWriter(name)
	if err != nil {
		return fmt.Errorf("criar stream: %w", err)
	}
	styles := newStyleCache(f, t)

	cols := s.ColCount()
	for c := range cols {
		w := t.ColWidth(c)
		if c < len(s.ColWidths) && s.ColWidths[c] > 0 {
			w = s.ColWidths[c]
		}
		if err := sw.SetColWidth(c+1, c+1, w); err != nil {
			return fmt.Errorf("largura da coluna: %w", err)
		}
	}
	if t.FreezeRow > 0 || t.FreezeCol > 0 {
		if err := sw.SetPanes(freezePanes(t)); err != nil {
			return fmt.Errorf("congelar painéis: %w", err)
		}
	}

	heights := map[int]float64{}
	for _, h := range t.RowHeights {
		heights[h.Row] = h.Height
	}
	total := s.RowCount()
	for r := range total {
		if r%progressStep == 0 {
			if err := ctx.Err(); err != nil {
				return err
			}
			if progress != nil {
				progress(r, total)
			}
		}
		row := make([]any, cols)
		for c := range cols {
			role := t.Role(r, c)
			cell := excelize.Cell{StyleID: styles.id(role)}
			v := s.Get(r, c)
			if formula.IsFormula(v) {
				cell.Formula = formula.ToExcel(v)
			} else if v != "" {
				cell.Value = typedValue(v, role)
			}
			row[c] = cell
		}
		h := t.DefaultRowHeight
		if hh, ok := heights[r]; ok {
			h = hh
		}
		opts := []excelize.RowOpts{}
		if h > 0 {
			opts = append(opts, excelize.RowOpts{Height: h})
		}
		if err := sw.SetRow(sheet.CellName(r, 0), row, opts...); err != nil {
			return fmt.Errorf("gravar linha %d: %w", r+1, err)
		}
	}
	if err := sw.Flush(); err != nil {
		return fmt.Errorf("finalizar planilha: %w", err)
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if err := writeAtomic(path, func(w io.Writer) error {
		_, err := f.WriteTo(w)
		return err
	}); err != nil {
		return fmt.Errorf("salvar arquivo: %w", err)
	}
	if progress != nil {
		progress(total, total)
	}
	return nil
}

func typedValue(v string, role templates.Role) any {
	switch role.Kind() {
	case templates.KindInteger:
		if n, ok := sheet.ParseInteger(v); ok {
			return n
		}
	case templates.KindDecimal, templates.KindCurrency:
		if n, ok := sheet.ParseNumber(v); ok {
			return n
		}
	case templates.KindText:
		return v
	}
	if role == templates.RoleEmpty {
		if n, ok := sheet.ParseNumber(v); ok {
			return n
		}
	}
	return v
}

func freezePanes(t *templates.Template) *excelize.Panes {
	return &excelize.Panes{
		Freeze:      true,
		XSplit:      t.FreezeCol,
		YSplit:      t.FreezeRow,
		TopLeftCell: sheet.CellName(t.FreezeRow, t.FreezeCol),
		ActivePane:  "bottomRight",
	}
}
