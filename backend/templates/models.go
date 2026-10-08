package templates

const currencyFmt = `"R$" #,##0.00`

func matrixStyles(extra map[string]Style) map[string]Style {
	styles := map[string]Style{
		"default":              {Size: 11},
		string(RoleClientName): {Size: 11, Wrap: true, VAlign: "bottom"},
		string(RoleTotal):      {Size: 11, HAlign: "center"},
	}
	for k, v := range extra {
		styles[k] = v
	}
	return styles
}

func preco() *Template {
	return &Template{
		ID:               "preco",
		Label:            "Preço",
		Description:      "Tabela de preços por cliente: produtos nas linhas, clientes nas colunas.",
		Layout:           Matrix,
		DefaultSheetName: "Planilha1",
		DefaultFileName:  "preco",
		FirstDataRow:     3,
		FirstDataCol:     2,
		HeaderRowRoles:   []Role{RoleRoute, RoleClientID, RoleClientName},
		Columns: []Column{
			{Header: "ID", Role: RoleProductID, Width: 9.14},
			{Header: "Produto", Role: RoleProductDesc, Width: 12.29},
		},
		DataRole:  RolePrice,
		DataWidth: 9.14,
		TotalCol:  -1,
		Fixed: []FixedCell{
			{Row: 0, Col: 0, Text: "PRECO", Role: RoleTitle, ReadOnly: true},
		},
		DefaultRowHeight: 15,
		RowHeights:       []RowHeight{{Row: 2, Height: 94.5}},
		Styles: matrixStyles(map[string]Style{
			string(RolePrice): {Size: 11, NumFmt: "#,##0.00"},
		}),
		DefaultClients:  13,
		DefaultProducts: 50,
	}
}

func venda() *Template {
	return &Template{
		ID:               "venda",
		Label:            "Venda",
		Description:      "Quantidades vendidas por cliente, com rota e soma por produto.",
		Layout:           Matrix,
		DefaultSheetName: "Planilha1",
		DefaultFileName:  "venda",
		FirstDataRow:     3,
		FirstDataCol:     3,
		HeaderRowRoles:   []Role{RoleRoute, RoleClientID, RoleClientName},
		Columns: []Column{
			{Header: "ID", Role: RoleProductID, Width: 9.14},
			{Header: "Produto", Role: RoleProductDesc, Width: 12.29},
			{Header: "SOMA", Role: RoleTotal, Width: 7.57},
		},
		DataRole:  RoleQuantity,
		DataWidth: 9.14,
		TotalCol:  2,
		Fixed: []FixedCell{
			{Row: 0, Col: 0, Text: "1", Role: RoleCode},
			{Row: 2, Col: 2, Text: "SOMA", Role: RoleHeader, ReadOnly: true},
		},
		DefaultRowHeight: 15,
		RowHeights:       []RowHeight{{Row: 2, Height: 94.5}},
		Styles: matrixStyles(map[string]Style{
			string(RoleHeader): {Size: 11, HAlign: "center"},
		}),
		DefaultClients:  13,
		DefaultProducts: 50,
	}
}

func compra() *Template {
	data := Style{Bold: true, Size: 11, HAlign: "center", Border: true}
	price := data
	price.NumFmt = currencyFmt
	return &Template{
		ID:               "compra",
		Label:            "Compra",
		Description:      "Pedido de compra: fornecedor, produto, quantidade e preço por linha.",
		Layout:           List,
		DefaultSheetName: "modelo pedido roça",
		DefaultFileName:  "compra",
		FirstDataRow:     2,
		Columns: []Column{
			{Header: "ID FORN", Role: RoleSupplierID, Width: 9},
			{Header: "FORNECEDOR", Role: RoleSupplierName, Width: 22.86},
			{Header: "ID PROD", Role: RoleProductID, Width: 9.14},
			{Header: "PRODUTO", Role: RoleProductDesc, Width: 18.14},
			{Header: "QTD", Role: RoleQuantity, Width: 14.29},
			{Header: "PRECO", Role: RolePrice, Width: 8.14},
		},
		TotalCol: -1,
		Fixed: []FixedCell{
			{Row: 0, Col: 0, Text: "COMPRA", Role: RoleTitle, ReadOnly: true},
			{Row: 1, Col: 0, Text: "ID FORN", Role: RoleHeader, ReadOnly: true},
			{Row: 1, Col: 1, Text: "FORNECEDOR", Role: RoleHeader, ReadOnly: true},
			{Row: 1, Col: 2, Text: "ID PROD", Role: RoleHeader, ReadOnly: true},
			{Row: 1, Col: 3, Text: "PRODUTO", Role: RoleHeader, ReadOnly: true},
			{Row: 1, Col: 4, Text: "QTD", Role: RoleHeader, ReadOnly: true},
			{Row: 1, Col: 5, Text: "PRECO", Role: RoleHeader, ReadOnly: true},
		},
		FreezeRow:        2,
		FreezeCol:        1,
		DefaultRowHeight: 18,
		Styles: map[string]Style{
			"default":                {Size: 11},
			string(RoleTitle):        {Bold: true, Size: 8},
			string(RoleHeader):       {Bold: true, Size: 12, HAlign: "center", VAlign: "center", Border: true},
			string(RoleSupplierID):   data,
			string(RoleSupplierName): data,
			string(RoleProductID):    data,
			string(RoleProductDesc):  data,
			string(RoleQuantity):     data,
			string(RolePrice):        price,
		},
		DefaultProducts: 150,
	}
}
