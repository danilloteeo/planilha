package excel

import (
	"context"
	"os"
	"path/filepath"
	"strconv"
	"testing"

	"github.com/xuri/excelize/v2"

	"planly/backend/sheet"
	"planly/backend/templates"
)

func samplePath(name string) string {
	return filepath.Join("..", "..", "docs", "models", name)
}

func TestImportSamples(t *testing.T) {
	cases := []struct {
		file, template string
		check          func(t *testing.T, s *sheet.Sheet)
	}{
		{"Preco.xlsx", "preco", func(t *testing.T, s *sheet.Sheet) {
			expect(t, s, 0, 0, "PRECO")
			expect(t, s, 1, 2, "1")
			expect(t, s, 2, 2, "MAGE")
			expect(t, s, 2, 14, "PARADA ANGELICA")
			expect(t, s, 3, 1, "ALHO PORÓ")
		}},
		{"venda.xlsx", "venda", func(t *testing.T, s *sheet.Sheet) {
			expect(t, s, 0, 0, "1")
			expect(t, s, 2, 2, "SOMA")
			expect(t, s, 3, 2, "=SOMA(D4:P4)")
			expect(t, s, 4, 2, "=SOMA(D5:P5)")
			expect(t, s, 4, 1, "LOURO")
		}},
		{"COMPRA.xlsx", "compra", func(t *testing.T, s *sheet.Sheet) {
			expect(t, s, 0, 0, "COMPRA")
			expect(t, s, 1, 5, "PRECO")
		}},
	}
	for _, tc := range cases {
		t.Run(tc.file, func(t *testing.T) {
			s, err := Import(context.Background(), samplePath(tc.file), nil)
			if err != nil {
				t.Fatalf("import: %v", err)
			}
			if s.Template != tc.template {
				t.Fatalf("template = %q, want %q", s.Template, tc.template)
			}
			tc.check(t, s)
		})
	}
}

func TestExportRoundTrip(t *testing.T) {
	for _, tpl := range templates.All() {
		t.Run(tpl.ID, func(t *testing.T) {
			s := tpl.New(templates.NewOptions{Name: "Teste", Clients: 3, Products: 4})
			if tpl.Layout == templates.Matrix {
				s.Set(1, tpl.FirstDataCol, "10")
				s.Set(2, tpl.FirstDataCol, "MAGE")
				s.Set(tpl.FirstDataRow, 0, "7")
				s.Set(tpl.FirstDataRow, 1, "ALHO")
				s.Set(tpl.FirstDataRow, tpl.FirstDataCol, "25,90")
			} else {
				s.Set(2, 0, "3")
				s.Set(2, 5, "1.234,50")
			}
			path := filepath.Join(t.TempDir(), "out.xlsx")
			if err := ExportXLSX(context.Background(), s, tpl, path, nil); err != nil {
				t.Fatalf("export: %v", err)
			}
			f, err := excelize.OpenFile(path)
			if err != nil {
				t.Fatalf("reopen: %v", err)
			}
			defer f.Close()
			name := f.GetSheetList()[0]
			if name != tpl.DefaultSheetName {
				t.Errorf("sheet name = %q, want %q", name, tpl.DefaultSheetName)
			}
			if tpl.ID == "venda" {
				fx, _ := f.GetCellFormula(name, "C4")
				if fx != "SUM(D4:F4)" {
					t.Errorf("C4 formula = %q", fx)
				}
			}
			if tpl.Layout == templates.Matrix {
				v, _ := f.GetCellValue(name, sheet.CellName(tpl.FirstDataRow, tpl.FirstDataCol), excelize.Options{RawCellValue: true})
				if n, _ := strconv.ParseFloat(v, 64); n != 25.9 {
					t.Errorf("price raw = %q, want 25.9", v)
				}
			} else {
				v, _ := f.GetCellValue(name, "F3", excelize.Options{RawCellValue: true})
				if v != "1234.5" {
					t.Errorf("F3 raw = %q, want 1234.5", v)
				}
				panes, _ := f.GetPanes(name)
				if !panes.Freeze || panes.TopLeftCell != "B3" {
					t.Errorf("panes = %+v", panes)
				}
			}
			back, err := Import(context.Background(), path, nil)
			if err != nil {
				t.Fatalf("reimport: %v", err)
			}
			if back.Template != tpl.ID {
				t.Errorf("reimport template = %q", back.Template)
			}
		})
	}
}

func TestCSVRoundTrip(t *testing.T) {
	tpl, _ := templates.Get("compra")
	s := tpl.New(templates.NewOptions{Products: 2})
	s.Set(2, 1, "Fornecedor; com ponto e vírgula")
	path := filepath.Join(t.TempDir(), "out.csv")
	if err := ExportCSV(context.Background(), s, path, nil); err != nil {
		t.Fatal(err)
	}
	back, err := Import(context.Background(), path, nil)
	if err != nil {
		t.Fatal(err)
	}
	expect(t, back, 2, 1, "Fornecedor; com ponto e vírgula")
}

func TestImportRejects(t *testing.T) {
	dir := t.TempDir()
	bad := filepath.Join(dir, "bad.xlsx")
	os.WriteFile(bad, []byte("not a zip"), 0o644)
	if _, err := Import(context.Background(), bad, nil); err == nil {
		t.Error("corrupted file accepted")
	}
	unknown := filepath.Join(dir, "x.csv")
	os.WriteFile(unknown, []byte("a;b\n1;2\n"), 0o644)
	if _, err := Import(context.Background(), unknown, nil); err != errNoTemplate {
		t.Errorf("unknown layout err = %v", err)
	}
	if _, err := Import(context.Background(), filepath.Join(dir, "x.xls"), nil); err == nil {
		t.Error("missing .xls accepted")
	}
}

func TestExportCanceled(t *testing.T) {
	tpl, _ := templates.Get("compra")
	s := tpl.New(templates.NewOptions{Products: 10})
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	path := filepath.Join(t.TempDir(), "x.xlsx")
	if err := ExportXLSX(ctx, s, tpl, path, nil); err == nil {
		t.Fatal("canceled export succeeded")
	}
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("canceled export left a file behind")
	}
}

func TestExportLarge(t *testing.T) {
	if testing.Short() {
		t.Skip()
	}
	tpl, _ := templates.Get("venda")
	s := tpl.New(templates.NewOptions{Clients: 13, Products: 100000})
	for r := tpl.FirstDataRow; r < s.RowCount(); r++ {
		s.Set(r, 0, strconv.Itoa(r))
		s.Set(r, 3, "5")
	}
	path := filepath.Join(t.TempDir(), "big.xlsx")
	if err := ExportXLSX(context.Background(), s, tpl, path, nil); err != nil {
		t.Fatal(err)
	}
	back, err := Import(context.Background(), path, nil)
	if err != nil {
		t.Fatal(err)
	}
	if back.RowCount() != s.RowCount() {
		t.Errorf("rows = %d, want %d", back.RowCount(), s.RowCount())
	}
}

func expect(t *testing.T, s *sheet.Sheet, r, c int, want string) {
	t.Helper()
	if got := s.Get(r, c); got != want {
		t.Errorf("%s = %q, want %q", sheet.CellName(r, c), got, want)
	}
}
