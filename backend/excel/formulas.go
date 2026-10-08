package excel

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"io"
	"path"
	"strings"

	"github.com/xuri/excelize/v2"
)

type formulaRef struct {
	axis   string
	text   string
	shared bool
}

func readFormulas(f *excelize.File, filePath, name string) (map[[2]int]string, error) {
	zr, err := zip.OpenReader(filePath)
	if err != nil {
		return nil, err
	}
	defer zr.Close()
	part, err := sheetPart(&zr.Reader, name)
	if err != nil {
		return nil, err
	}
	refs, err := scanFormulas(&zr.Reader, part)
	if err != nil {
		return nil, err
	}
	out := make(map[[2]int]string, len(refs))
	for _, ref := range refs {
		col, row, err := excelize.CellNameToCoordinates(ref.axis)
		if err != nil {
			continue
		}
		text := ref.text
		if ref.shared && text == "" {
			if text, err = f.GetCellFormula(name, ref.axis); err != nil {
				continue
			}
		}
		if text != "" {
			out[[2]int{row - 1, col - 1}] = text
		}
	}
	return out, nil
}

func sheetPart(zr *zip.Reader, name string) (string, error) {
	var wb struct {
		Sheets []struct {
			Name string `xml:"name,attr"`
			RID  string `xml:"http://schemas.openxmlformats.org/officeDocument/2006/relationships id,attr"`
		} `xml:"sheets>sheet"`
	}
	if err := decodePart(zr, "xl/workbook.xml", &wb); err != nil {
		return "", err
	}
	var rels struct {
		Items []struct {
			ID     string `xml:"Id,attr"`
			Target string `xml:"Target,attr"`
		} `xml:"Relationship"`
	}
	if err := decodePart(zr, "xl/_rels/workbook.xml.rels", &rels); err != nil {
		return "", err
	}
	for _, s := range wb.Sheets {
		if s.Name != name {
			continue
		}
		for _, rel := range rels.Items {
			if rel.ID == s.RID {
				if strings.HasPrefix(rel.Target, "/") {
					return strings.TrimPrefix(rel.Target, "/"), nil
				}
				return path.Join("xl", rel.Target), nil
			}
		}
	}
	return "", errors.New("aba não encontrada no pacote")
}

func decodePart(zr *zip.Reader, name string, v any) error {
	rc, err := openPart(zr, name)
	if err != nil {
		return err
	}
	defer rc.Close()
	return xml.NewDecoder(rc).Decode(v)
}

func openPart(zr *zip.Reader, name string) (io.ReadCloser, error) {
	for _, file := range zr.File {
		if strings.EqualFold(file.Name, name) {
			return file.Open()
		}
	}
	return nil, errors.New("parte não encontrada: " + name)
}

func scanFormulas(zr *zip.Reader, part string) ([]formulaRef, error) {
	rc, err := openPart(zr, part)
	if err != nil {
		return nil, err
	}
	defer rc.Close()
	dec := xml.NewDecoder(rc)
	var refs []formulaRef
	axis := ""
	for {
		tok, err := dec.Token()
		if err == io.EOF {
			return refs, nil
		}
		if err != nil {
			return nil, err
		}
		el, ok := tok.(xml.StartElement)
		if !ok {
			continue
		}
		switch el.Name.Local {
		case "c":
			axis = attr(el, "r")
		case "f":
			var body string
			if err := dec.DecodeElement(&body, &el); err != nil {
				return nil, err
			}
			refs = append(refs, formulaRef{axis: axis, text: body, shared: attr(el, "t") == "shared"})
		}
	}
}

func attr(el xml.StartElement, name string) string {
	for _, a := range el.Attr {
		if a.Name.Local == name {
			return a.Value
		}
	}
	return ""
}
