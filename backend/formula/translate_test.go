package formula

import "testing"

func TestToExcel(t *testing.T) {
	cases := map[string]string{
		"=SOMA(D4:P4)":                   "SUM(D4:P4)",
		"=soma(D4:P4)":                   "SUM(D4:P4)",
		"=MEDIA(A1:A3)":                  "AVERAGE(A1:A3)",
		"=SE(A1>1,5;\"sim; não\";B2)":    "IF(A1>1.5,\"sim; não\",B2)",
		"=ARRED(B2*C2;2)":                "ROUND(B2*C2,2)",
		"=CONT.SE(A:A;\">10\")":          "COUNTIF(A:A,\">10\")",
		"=PROCV(A2;'Aba 1'!A:B;2;0)":     "VLOOKUP(A2,'Aba 1'!A:B,2,0)",
		"=$B$2*C2":                       "$B$2*C2",
		"=HOJE()":                        "TODAY()",
		"=MÁXIMO(D2:D20)-MÍNIMO(D2:D20)": "MAX(D2:D20)-MIN(D2:D20)",
	}
	for in, want := range cases {
		if got := ToExcel(in); got != want {
			t.Errorf("ToExcel(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestFromExcel(t *testing.T) {
	cases := map[string]string{
		"SUM(D4:P4)":                 "=SOMA(D4:P4)",
		"IF(A1>1.5,\"a,b\",B2)":      "=SE(A1>1,5;\"a,b\";B2)",
		"VLOOKUP(A2,Sheet2!A:B,2,0)": "=PROCV(A2;Sheet2!A:B;2;0)",
		"UNKNOWNFN(A1)":              "=UNKNOWNFN(A1)",
	}
	for in, want := range cases {
		if got := FromExcel(in); got != want {
			t.Errorf("FromExcel(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestRoundTrip(t *testing.T) {
	for _, f := range []string{"=SOMA(D4:P4)", "=SE(A1>2,5;1;0)", "=MÉDIA(A1:A9)*1,1"} {
		if got := FromExcel(ToExcel(f)); got != f {
			t.Errorf("round trip %q -> %q", f, got)
		}
	}
}

func TestCheck(t *testing.T) {
	valid := []string{"=SOMA(A1:A3)", "=A1+B1", "=SE(A1=\"x\";1;0)", "=1,5*2", "=\"a\"\"b\""}
	for _, v := range valid {
		if err := Check(v); err != nil {
			t.Errorf("Check(%q) = %v, want nil", v, err)
		}
	}
	invalid := []string{"=", "=SOMA(A1:A3", "=A1+", "=FOO(A1)", "=\"abc", "=)A1(", "=\"a\"\""}
	for _, v := range invalid {
		if err := Check(v); err == nil {
			t.Errorf("Check(%q) = nil, want error", v)
		}
	}
}
