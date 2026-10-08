package formula

import (
	"strings"
	"unicode"

	"golang.org/x/text/runes"
	"golang.org/x/text/transform"
	"golang.org/x/text/unicode/norm"
)

type Function struct {
	PT          string `json:"pt"`
	EN          string `json:"en"`
	Syntax      string `json:"syntax"`
	Description string `json:"description"`
}

var Functions = []Function{
	{"SOMA", "SUM", "SOMA(intervalo)", "Soma os valores do intervalo"},
	{"MÉDIA", "AVERAGE", "MÉDIA(intervalo)", "Média dos valores do intervalo"},
	{"MÍNIMO", "MIN", "MÍNIMO(intervalo)", "Menor valor do intervalo"},
	{"MÁXIMO", "MAX", "MÁXIMO(intervalo)", "Maior valor do intervalo"},
	{"CONT.SE", "COUNTIF", "CONT.SE(intervalo; critério)", "Conta as células que atendem ao critério"},
	{"SE", "IF", "SE(teste; se_verdadeiro; se_falso)", "Retorna um valor conforme a condição"},
	{"PROCV", "VLOOKUP", "PROCV(valor; tabela; coluna; 0)", "Procura um valor na primeira coluna da tabela"},
	{"SOMASE", "SUMIF", "SOMASE(intervalo; critério; intervalo_soma)", "Soma as células que atendem ao critério"},
	{"ARRED", "ROUND", "ARRED(número; casas)", "Arredonda o número"},
	{"HOJE", "TODAY", "HOJE()", "Data de hoje"},
	{"CONT.VALORES", "COUNTA", "CONT.VALORES(intervalo)", "Conta as células preenchidas"},
	{"CONT.NÚM", "COUNT", "CONT.NÚM(intervalo)", "Conta as células com números"},
	{"CONT.SES", "COUNTIFS", "CONT.SES(intervalo; critério; ...)", "Conta com vários critérios"},
	{"SOMASES", "SUMIFS", "SOMASES(intervalo_soma; intervalo; critério; ...)", "Soma com vários critérios"},
	{"MÉDIASE", "AVERAGEIF", "MÉDIASE(intervalo; critério)", "Média das células que atendem ao critério"},
	{"PROCH", "HLOOKUP", "PROCH(valor; tabela; linha; 0)", "Procura um valor na primeira linha da tabela"},
	{"SEERRO", "IFERROR", "SEERRO(valor; se_erro)", "Retorna outro valor em caso de erro"},
	{"E", "AND", "E(condição; ...)", "Verdadeiro se todas as condições forem verdadeiras"},
	{"OU", "OR", "OU(condição; ...)", "Verdadeiro se alguma condição for verdadeira"},
	{"NÃO", "NOT", "NÃO(condição)", "Inverte a condição"},
	{"ABS", "ABS", "ABS(número)", "Valor absoluto"},
	{"INT", "INT", "INT(número)", "Parte inteira"},
	{"MOD", "MOD", "MOD(número; divisor)", "Resto da divisão"},
	{"POTÊNCIA", "POWER", "POTÊNCIA(número; expoente)", "Potência"},
	{"RAIZ", "SQRT", "RAIZ(número)", "Raiz quadrada"},
	{"MULT", "PRODUCT", "MULT(intervalo)", "Multiplica os valores"},
	{"ARREDONDAR.PARA.CIMA", "ROUNDUP", "ARREDONDAR.PARA.CIMA(número; casas)", "Arredonda para cima"},
	{"ARREDONDAR.PARA.BAIXO", "ROUNDDOWN", "ARREDONDAR.PARA.BAIXO(número; casas)", "Arredonda para baixo"},
	{"AGORA", "NOW", "AGORA()", "Data e hora atuais"},
	{"DATA", "DATE", "DATA(ano; mês; dia)", "Monta uma data"},
	{"ANO", "YEAR", "ANO(data)", "Ano da data"},
	{"MÊS", "MONTH", "MÊS(data)", "Mês da data"},
	{"DIA", "DAY", "DIA(data)", "Dia da data"},
	{"CONCATENAR", "CONCATENATE", "CONCATENAR(texto; ...)", "Junta textos"},
	{"ESQUERDA", "LEFT", "ESQUERDA(texto; n)", "Primeiros caracteres"},
	{"DIREITA", "RIGHT", "DIREITA(texto; n)", "Últimos caracteres"},
	{"EXT.TEXTO", "MID", "EXT.TEXTO(texto; início; n)", "Trecho do texto"},
	{"MAIÚSCULA", "UPPER", "MAIÚSCULA(texto)", "Converte para maiúsculas"},
	{"MINÚSCULA", "LOWER", "MINÚSCULA(texto)", "Converte para minúsculas"},
	{"ARRUMAR", "TRIM", "ARRUMAR(texto)", "Remove espaços extras"},
	{"NÚM.CARACT", "LEN", "NÚM.CARACT(texto)", "Quantidade de caracteres"},
	{"TEXTO", "TEXT", "TEXTO(valor; formato)", "Formata um valor como texto"},
	{"VALOR", "VALUE", "VALOR(texto)", "Converte texto em número"},
	{"ÍNDICE", "INDEX", "ÍNDICE(tabela; linha; coluna)", "Valor em uma posição da tabela"},
	{"CORRESP", "MATCH", "CORRESP(valor; intervalo; 0)", "Posição de um valor no intervalo"},
}

var (
	ptToEN = map[string]string{}
	enToPT = map[string]string{}
)

func init() {
	for _, f := range Functions {
		ptToEN[foldName(f.PT)] = f.EN
		enToPT[f.EN] = f.PT
	}
}

func foldName(s string) string {
	t := transform.Chain(norm.NFD, runes.Remove(runes.In(unicode.Mn)), norm.NFC)
	out, _, err := transform.String(t, strings.ToUpper(s))
	if err != nil {
		return strings.ToUpper(s)
	}
	return out
}

func LookupPT(name string) (string, bool) {
	en, ok := ptToEN[foldName(name)]
	return en, ok
}

func LookupEN(name string) (string, bool) {
	pt, ok := enToPT[strings.ToUpper(name)]
	return pt, ok
}
