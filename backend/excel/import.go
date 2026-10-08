package excel

import (
	"bytes"
	"context"
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/xuri/excelize/v2"

	"planly/backend/apperr"
	"planly/backend/formula"
	"planly/backend/sheet"
	"planly/backend/templates"
)

const MaxImportSize = 64 << 20

var errNoTemplate = apperr.New("O arquivo não segue nenhum dos modelos (Preço, Venda ou Compra).", nil)

func Import(ctx context.Context, path string, progress Progress) (*sheet.Sheet, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, apperr.New("Não foi possível acessar o arquivo selecionado.", err)
	}
	if info.Size() > MaxImportSize {
		return nil, apperr.New("O arquivo é grande demais para importar (limite de 64 MB).", nil)
	}
	switch strings.ToLower(filepath.Ext(path)) {
	case ".xlsx":
		return importXLSX(ctx, path, progress)
	case ".csv":
		return importCSV(path)
	default:
		return nil, apperr.New("Formato não suportado. Use arquivos .xlsx ou .csv.", nil)
	}
}

func importXLSX(ctx context.Context, path string, progress Progress) (*sheet.Sheet, error) {
	f, err := excelize.OpenFile(path)
	if err != nil {
		return nil, apperr.New("O arquivo está corrompido ou não é uma planilha Excel válida.", err)
	}
	defer f.Close()

	sheets := f.GetSheetList()
	if len(sheets) == 0 {
		return nil, apperr.New("A planilha não possui abas.", nil)
	}
	name := sheets[0]
	rows, err := f.GetRows(name, excelize.Options{RawCellValue: true})
	if err != nil {
		return nil, apperr.New("Não foi possível ler os dados da planilha.", err)
	}
	t, ok := templates.Detect(rows)
	if !ok {
		return nil, errNoTemplate
	}
	formulas, err := readFormulas(f, path, name)
	if err != nil {
		return nil, apperr.New("Não foi possível ler as fórmulas da planilha.", err)
	}

	s := &sheet.Sheet{Template: t.ID, Name: strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)), SheetName: name}
	total := len(rows)
	s.Cells = make([][]string, total)
	for r, row := range rows {
		if r%progressStep == 0 {
			if err := ctx.Err(); err != nil {
				return nil, err
			}
			if progress != nil {
				progress(r, total)
			}
		}
		out := make([]string, len(row))
		for c, v := range row {
			if fx, ok := formulas[[2]int{r, c}]; ok {
				out[c] = formula.FromExcel(fx)
				continue
			}
			out[c] = normalizeValue(v, t.Role(r, c))
		}
		s.Cells[r] = out
	}
	for pos, fx := range formulas {
		if pos[0] >= len(s.Cells) || pos[1] >= len(s.Cells[pos[0]]) {
			s.Set(pos[0], pos[1], formula.FromExcel(fx))
		}
	}
	readColWidths(f, name, s)
	t.Normalize(s)
	if progress != nil {
		progress(total, total)
	}
	return s, nil
}

func readColWidths(f *excelize.File, name string, s *sheet.Sheet) {
	cols := s.ColCount()
	s.ColWidths = make([]float64, cols)
	for c := range cols {
		if w, err := f.GetColWidth(name, sheet.ColName(c)); err == nil {
			s.ColWidths[c] = w
		}
	}
}

func normalizeValue(v string, role templates.Role) string {
	switch role.Kind() {
	case templates.KindInteger, templates.KindDecimal, templates.KindCurrency:
		if out, ok := sheet.FromRawNumber(v); ok {
			return out
		}
	}
	return v
}

func importCSV(path string) (*sheet.Sheet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, apperr.New("Não foi possível ler o arquivo CSV.", err)
	}
	data = bytes.TrimPrefix(data, []byte("\xef\xbb\xbf"))
	if !utf8.Valid(data) {
		data = latin1ToUTF8(data)
	}
	r := csv.NewReader(bytes.NewReader(data))
	r.Comma = sniffDelimiter(data)
	r.FieldsPerRecord = -1
	r.LazyQuotes = true
	rows, err := r.ReadAll()
	if err != nil {
		var pe *csv.ParseError
		if errors.As(err, &pe) {
			return nil, apperr.New(fmt.Sprintf("Erro ao ler o CSV na linha %d.", pe.Line), err)
		}
		return nil, apperr.New("O arquivo CSV está corrompido.", err)
	}
	t, ok := templates.Detect(rows)
	if !ok {
		return nil, errNoTemplate
	}
	s := &sheet.Sheet{
		Template:  t.ID,
		Name:      strings.TrimSuffix(filepath.Base(path), filepath.Ext(path)),
		SheetName: t.DefaultSheetName,
		Cells:     rows,
	}
	t.Normalize(s)
	return s, nil
}

func sniffDelimiter(data []byte) rune {
	line, _, _ := bytes.Cut(data, []byte("\n"))
	best, bestCount := ';', -1
	for _, d := range []rune{';', ',', '\t'} {
		if n := bytes.Count(line, []byte(string(d))); n > bestCount {
			best, bestCount = d, n
		}
	}
	return best
}

func latin1ToUTF8(data []byte) []byte {
	var b bytes.Buffer
	for _, c := range data {
		b.WriteRune(rune(c))
	}
	return b.Bytes()
}

func ExportCSV(ctx context.Context, s *sheet.Sheet, path string, progress Progress) error {
	return writeAtomic(path, func(w io.Writer) error {
		if _, err := w.Write([]byte("\xef\xbb\xbf")); err != nil {
			return err
		}
		cw := csv.NewWriter(w)
		cw.Comma = ';'
		cols := s.ColCount()
		total := s.RowCount()
		record := make([]string, cols)
		for r := range total {
			if r%progressStep == 0 {
				if err := ctx.Err(); err != nil {
					return err
				}
				if progress != nil {
					progress(r, total)
				}
			}
			for c := range cols {
				record[c] = s.Get(r, c)
			}
			if err := cw.Write(record); err != nil {
				return err
			}
		}
		cw.Flush()
		if err := cw.Error(); err != nil {
			return err
		}
		if progress != nil {
			progress(total, total)
		}
		return nil
	})
}
