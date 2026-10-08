package excel

import (
	"context"
	"path/filepath"
	"strconv"
	"testing"
	"time"

	"planly/backend/templates"
)

func TestPerfBreakdown(t *testing.T) {
	tpl, _ := templates.Get("venda")
	s := tpl.New(templates.NewOptions{Clients: 13, Products: 100000})
	for r := tpl.FirstDataRow; r < s.RowCount(); r++ {
		s.Set(r, 0, strconv.Itoa(r))
		s.Set(r, 3, "5")
	}
	path := filepath.Join(t.TempDir(), "big.xlsx")
	t0 := time.Now()
	ExportXLSX(context.Background(), s, tpl, path, nil)
	t.Logf("export %v", time.Since(t0))
	t0 = time.Now()
	Import(context.Background(), path, nil)
	t.Logf("import %v", time.Since(t0))
}
