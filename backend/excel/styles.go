package excel

import (
	"github.com/xuri/excelize/v2"

	"planly/backend/templates"
)

type styleCache struct {
	f   *excelize.File
	t   *templates.Template
	ids map[templates.Role]int
}

func newStyleCache(f *excelize.File, t *templates.Template) *styleCache {
	return &styleCache{f: f, t: t, ids: map[templates.Role]int{}}
}

func (c *styleCache) id(role templates.Role) int {
	if id, ok := c.ids[role]; ok {
		return id
	}
	id, err := c.f.NewStyle(toExcelStyle(c.t.StyleFor(role)))
	if err != nil {
		id = 0
	}
	c.ids[role] = id
	return id
}

func toExcelStyle(st templates.Style) *excelize.Style {
	size := st.Size
	if size == 0 {
		size = 11
	}
	out := &excelize.Style{
		Font: &excelize.Font{Family: "Calibri", Size: size, Bold: st.Bold},
	}
	if st.HAlign != "" || st.VAlign != "" || st.Wrap {
		out.Alignment = &excelize.Alignment{
			Horizontal: st.HAlign,
			Vertical:   st.VAlign,
			WrapText:   st.Wrap,
		}
	}
	if st.Border {
		for _, side := range []string{"left", "top", "right", "bottom"} {
			out.Border = append(out.Border, excelize.Border{Type: side, Color: "000000", Style: 1})
		}
	}
	if st.NumFmt != "" {
		numFmt := st.NumFmt
		out.CustomNumFmt = &numFmt
	}
	return out
}
