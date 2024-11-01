package arvore

import (
	"ChikenInterpreter/ferramentas"
	lex "ChikenInterpreter/lexing"
	"fmt"
	"strings"
)

// Expressoes compostas:
// expressoes de prefixo
type ExpressaodePrefixo struct {
	Token      lex.Token
	Operador   string
	ExpDireita Expressao
	Profu      int
}

func (op *ExpressaodePrefixo) noExpressao() {}
func (op *ExpressaodePrefixo) IncProfu() {
	op.Profu++
	op.ExpDireita.IncProfu()
}
func (op *ExpressaodePrefixo) GetTokenNo() lex.Token {
	return op.Token
}
func (op *ExpressaodePrefixo) GetInformacao() string {
	identacao := ferramentas.GetIdentacao(op.Profu)

	texto := identacao
	texto += "Operador prefixo: " + op.Operador
	texto += "\n" + identacao + "Expressao:\n"
	texto += op.ExpDireita.GetInformacao()
	return texto
}

// expressoes de infixo
type ExpressaoInfixo struct {
	Token       lex.Token
	Operador    string
	ExpDireita  Expressao
	ExpEsquerda Expressao
	Profu       int
}

func (ei *ExpressaoInfixo) noExpressao() {}
func (ei *ExpressaoInfixo) IncProfu() {
	ei.Profu++
	ei.ExpEsquerda.IncProfu()
	ei.ExpDireita.IncProfu()
}
func (ei *ExpressaoInfixo) GetTokenNo() lex.Token {
	return ei.Token
}
func (ei *ExpressaoInfixo) GetInformacao() string {
	identacao := ferramentas.GetIdentacao(ei.Profu)

	texto := identacao
	texto += "Operador infixo: " + ei.Operador + "\n"
	texto += identacao + "Expressao esquerda:\n"
	texto += ei.ExpEsquerda.GetInformacao() + "\n"
	texto += identacao + "Expressao direita:\n"
	texto += ei.ExpDireita.GetInformacao()

	return texto
}

type ExpressaoAtributo struct {
	Token    lex.Token
	Expres   Expressao
	Atributo string
	Profu    int
}

func (ea *ExpressaoAtributo) noExpressao() {}
func (ea *ExpressaoAtributo) IncProfu() {
	ea.Expres.IncProfu()
	ea.Profu++
}
func (ea *ExpressaoAtributo) GetTokenNo() lex.Token {
	return ea.Token
}
func (ea *ExpressaoAtributo) GetInformacao() string {
	parts := make([]string, 3)
	indent := ferramentas.GetIdentacao(ea.Profu)
	parts[0] = indent + "Expressao de captura de atributo. Expressao:"
	parts[1] = ea.Expres.GetInformacao()
	parts[2] = indent + "Atributo: " + ea.Atributo

	return strings.Join(parts, "\n")
}

// Expressos com blocos:
// Expressao if
type ExpressaoIf struct {
	Token      lex.Token
	Profu      int
	Condicao   Expressao
	BlocoEntao *BlocoInstrucao
	BlocoSenao *BlocoInstrucao
}

func (ei *ExpressaoIf) noExpressao() {}
func (ei *ExpressaoIf) IncProfu() {
	ei.Condicao.IncProfu()
	ei.BlocoEntao.IncProfu()
	if ei.BlocoSenao != nil {
		ei.BlocoSenao.IncProfu()
	}
}
func (ei *ExpressaoIf) GetTokenNo() lex.Token {
	return ei.Token
}
func (ei *ExpressaoIf) GetInformacao() string {
	indent := ferramentas.GetIdentacao(ei.Profu)
	var parts []string = make([]string, 7)

	parts[0] = indent + "Expressao IF:"
	parts[1] = indent + "Condicao:"
	parts[2] = ei.Condicao.GetInformacao()
	parts[3] = indent + "Bloco entao:"
	parts[4] = ei.BlocoEntao.GetInformacao()
	parts[5] = indent + "Bloco senao:"
	if ei.BlocoSenao != nil {
		parts[6] = ei.BlocoSenao.GetInformacao()
	} else {
		parts[6] = indent + "    Bloco vazio"
	}

	return strings.Join(parts, "\n")
}

// Expressao repeat
type ExpressaoRepeat struct {
	Token        lex.Token
	Condicao1    Expressao
	Condicao2    Expressao
	BlocoRepetir *BlocoInstrucao
	Profu        int
}

func (er *ExpressaoRepeat) noExpressao() {}
func (er *ExpressaoRepeat) IncProfu() {
	if er.Condicao1 != nil {
		er.Condicao1.IncProfu()
	}
	if er.Condicao2 != nil {
		er.Condicao2.IncProfu()
	}
	er.BlocoRepetir.IncProfu()
}
func (er *ExpressaoRepeat) GetTokenNo() lex.Token {
	return er.Token
}
func (er *ExpressaoRepeat) GetInformacao() string {
	indent := ferramentas.GetIdentacao(er.Profu)
	parts := make([]string, 4)

	parts[0] = indent + "Expressao REPEAT."
	if er.Condicao1 != nil {
		parts[1] = indent + "Condicao 1:\n" + er.Condicao1.GetInformacao()
	} else {
		parts[1] = indent + "Sem a primeira condicao."
	}

	if er.Condicao2 != nil {
		parts[2] = indent + "Condicao 2:\n" + er.Condicao2.GetInformacao()
	} else {
		parts[2] = indent + "Sem a segunda condicao."
	}

	parts[3] = indent + "Bloco de instrucoes:\n" + er.BlocoRepetir.GetInformacao()

	return strings.Join(parts, "\n")
}

// Expressoes fun
type ExpressaoFun struct {
	Token      lex.Token
	Parametros []*Identificador
	Bloco      *BlocoInstrucao
	Profu      int
}

func (ef *ExpressaoFun) noExpressao() {}
func (ef *ExpressaoFun) IncProfu() {
	ef.Profu++
	ef.Bloco.IncProfu()
	for _, i := range ef.Parametros {
		i.IncProfu()
	}
}
func (ef *ExpressaoFun) GetTokenNo() lex.Token {
	return ef.Token
}
func (ef *ExpressaoFun) GetInformacao() string {
	indent := ferramentas.GetIdentacao(ef.Profu)
	var parts []string = make([]string, 5)
	var parametros string

	for _, i := range ef.Parametros {
		parametros += i.GetInformacao() + ",\n"
	}

	parts[0] = indent + "Expressao FUN"
	parts[1] = indent + "Parametros:"
	parts[2] = parametros
	parts[3] = indent + "Bloco de instrucoes:"
	parts[4] = ef.Bloco.GetInformacao()

	return strings.Join(parts, "\n")
}

// Chamada de funcao
type CallFun struct {
	Token      lex.Token
	Funcao     Expressao
	Argumentos []Expressao
	Profu      int
}

func (cf *CallFun) noExpressao() {}
func (cf *CallFun) IncProfu() {
	if len(cf.Argumentos) > 0 {
		for _, i := range cf.Argumentos {
			i.IncProfu()
		}
	}

	cf.Funcao.IncProfu()
	cf.Profu++
}

func (cf *CallFun) GetTokenNo() lex.Token {
	return cf.Token
}
func (cf *CallFun) GetInformacao() string {
	var argumentos []string
	var texto string
	indent := ferramentas.GetIdentacao(cf.Profu)

	if len(cf.Argumentos) > 0 {
		for _, i := range cf.Argumentos {
			argumentos = append(argumentos, i.GetInformacao())
		}
	}

	texto = indent + "Chamada de funcao:\n" + cf.Funcao.GetInformacao()
	texto += "\n" + indent + "Argumentos:\n"
	texto += strings.Join(argumentos, ",\n")

	return texto
}

// Expressoes de dados complexos
// Listas:
type ExpressaoLista struct {
	Token      lex.Token
	Expressoes []Expressao
	Profu      int
}

func (el *ExpressaoLista) noExpressao() {}
func (el *ExpressaoLista) IncProfu() {
	el.Profu++
	for _, e := range el.Expressoes {
		e.IncProfu()
	}
}
func (el *ExpressaoLista) GetTokenNo() lex.Token {
	return el.Token
}
func (el *ExpressaoLista) GetInformacao() string {
	parts := make([]string, len(el.Expressoes)+1)
	indent := ferramentas.GetIdentacao(el.Profu)

	parts[0] = indent + "Expressao de lista. Valores:"

	for i, exp := range el.Expressoes {
		parts[i+1] = exp.GetInformacao()
	}

	return strings.Join(parts, "\n\n")
}

// Expressao de dicionario
type ExpressaoDict struct {
	Token   lex.Token
	Chaves  []Expressao
	Valores []Expressao
	Profu   int
}

func (em *ExpressaoDict) noExpressao() {}
func (em *ExpressaoDict) IncProfu() {
	for i, obj := range em.Chaves {
		obj.IncProfu()
		em.Valores[i].IncProfu()
	}

	em.Profu++
}
func (em *ExpressaoDict) GetTokenNo() lex.Token { return em.Token }
func (em *ExpressaoDict) GetInformacao() string {
	indent := ferramentas.GetIdentacao(em.Profu)
	parts := make([]string, len(em.Chaves)*2+1)

	parts[0] = indent + "Expressao de mapa de dados. Valores:"
	for i, chave := range em.Chaves {
		pos := 1 + (i * 2)
		parts[pos] = indent + "Chave:\n" + chave.GetInformacao()
		parts[pos+1] = indent + "Valor:\n" + em.Valores[i].GetInformacao()
	}

	return strings.Join(parts, "\n")
}

// Expressao de objeto:
type ExpressaoObjeto struct {
	Token     lex.Token
	Atributos map[string]Expressao
	Profu     int
}

func (eo *ExpressaoObjeto) noExpressao()          {}
func (eo *ExpressaoObjeto) IncProfu()             { eo.Profu++ }
func (eo *ExpressaoObjeto) GetTokenNo() lex.Token { return eo.Token }
func (eo *ExpressaoObjeto) GetInformacao() string {
	indent := ferramentas.GetIdentacao(eo.Profu)
	parts := make([]string, 1+len(eo.Atributos))
	i := 1

	parts[0] = indent + "Expressao OBJECT. Propriedades:"

	for chave, valor := range eo.Atributos {
		parts[i] = fmt.Sprintf("%s'%s':\n%s", indent, chave, valor.GetInformacao())
		i++
	}

	return strings.Join(parts, "\n")
}

type ExpressaoClass struct {
	Token          lex.Token
	Objeto         *ExpressaoObjeto
	SuperClasses   []Expressao
	AtributosObj   map[string]Expressao
	AtributosClass map[string]Expressao

	Profu int
}

func (ec *ExpressaoClass) SetAtribObj(nome string, valor Expressao) {
	ec.AtributosObj[nome] = valor
}
func (ec *ExpressaoClass) SetAtribClass(nome string, valor Expressao) {
	ec.AtributosClass[nome] = valor
}
func (ec *ExpressaoClass) noExpressao() {}
func (ec *ExpressaoClass) IncProfu() {
	ec.Objeto.IncProfu()

	for _, expr := range ec.AtributosObj {
		expr.IncProfu()
	}

	for _, expr := range ec.AtributosClass {
		expr.IncProfu()
	}

	ec.Profu++
}
func (ec *ExpressaoClass) GetTokenNo() lex.Token {
	return ec.Token
}

func (ec *ExpressaoClass) GetInformacao() string {
	indent := ferramentas.GetIdentacao(ec.Profu)
	parts := make([]string, len(ec.AtributosClass)+len(ec.AtributosObj)+6)
	var i int = 5

	parts[0] = indent + "Expressao de classe."

	if ec.Objeto != nil {
		parts[1] = indent + "    Objeto:\n" + ec.Objeto.GetInformacao()
	}

	parts[2] = fmt.Sprintf(indent+"    Número de SuperClases: %d", len(ec.SuperClasses))

	parts[4] = indent + "Atributos do Objeto:"

	for nome, valor := range ec.AtributosObj {
		parts[i] = fmt.Sprintf(indent+"    '%s'. Expressao:\n%s", nome, valor.GetInformacao())
		i++
	}

	parts[i] = indent + "Atributos de classe:"
	i++

	for nome, valor := range ec.AtributosClass {
		parts[i] = fmt.Sprintf(indent+"    '%s'. Expressao:\n%s", nome, valor.GetInformacao())
		i++
	}

	return strings.Join(parts, "\n")
}

func NewExpressaoClass(token lex.Token, profu int) *ExpressaoClass {
	novaExpre := &ExpressaoClass{Token: token, Profu: profu}
	novaExpre.AtributosObj = make(map[string]Expressao)
	novaExpre.AtributosClass = make(map[string]Expressao)

	return novaExpre
}

// Expressoes simples:
// instrucao de identificador
type ChamadaObjeto struct {
	Token lex.Token
	Profu int
}

func (co *ChamadaObjeto) noExpressao()          {}
func (co *ChamadaObjeto) IncProfu()             { co.Profu++ }
func (co *ChamadaObjeto) GetTokenNo() lex.Token { return co.Token }
func (co *ChamadaObjeto) GetInformacao() string {
	return ferramentas.GetIdentacao(co.Profu) + "Chamada de objeto."
}

type Identificador struct {
	Token lex.Token
	Nome  string
	Profu int
}

func (i *Identificador) noExpressao() {}
func (i *Identificador) IncProfu() {
	i.Profu++
}
func (i *Identificador) GetTokenNo() lex.Token {
	return i.Token
}
func (i *Identificador) GetInformacao() string {
	return ferramentas.GetIdentacao(i.Profu) + "Nome:" + i.Nome
}

// instrucao de real literal
type LiteralReal struct {
	Token lex.Token
	Valor float64
	Profu int
}

func (lr *LiteralReal) noExpressao() {}
func (lr *LiteralReal) IncProfu() {
	lr.Profu++
}
func (lr *LiteralReal) GetTokenNo() lex.Token {
	return lr.Token
}
func (lr *LiteralReal) GetInformacao() string {
	return ferramentas.GetIdentacao(lr.Profu) + "Valor: " + lr.Token.Valor
}

// instrucao de inteiro literal
type LiteralInt struct {
	Token lex.Token
	Valor int64
	Profu int
}

func (li *LiteralInt) noExpressao() {}
func (li *LiteralInt) IncProfu() {
	li.Profu++
}
func (li *LiteralInt) GetTokenNo() lex.Token {
	return li.Token
}
func (li *LiteralInt) GetInformacao() string {
	return ferramentas.GetIdentacao(li.Profu) + "Valor: " + li.Token.Valor
}

// instrucao booleana
type Booleano struct {
	Token lex.Token
	Valor bool
	Profu int
}

func (b *Booleano) noExpressao() {}
func (b *Booleano) IncProfu() {
	b.Profu++
}
func (b *Booleano) GetTokenNo() lex.Token {
	return b.Token
}
func (b *Booleano) GetInformacao() string {
	valor := ""
	if b.Token.Tipo == lex.TRUE {
		valor = "true"
	} else if b.Token.Tipo == lex.FALSE {
		valor = "false"
	}

	return ferramentas.GetIdentacao(b.Profu) + fmt.Sprintf("Token: %s, valor: %s", b.Token.Tipo, valor)
}

type LiteralString struct {
	Token lex.Token
	Valor string
	Profu int
}

func (ls *LiteralString) noExpressao()          {}
func (ls *LiteralString) IncProfu()             { ls.Profu++ }
func (ls *LiteralString) GetTokenNo() lex.Token { return ls.Token }
func (ls *LiteralString) GetInformacao() string {
	indent := ferramentas.GetIdentacao(ls.Profu)
	return indent + "String: " + ls.Valor
}

type TipoNone struct {
	Token lex.Token
	Profu int
}

func (tn *TipoNone) noExpressao()          {}
func (tn *TipoNone) IncProfu()             { tn.Profu++ }
func (tn *TipoNone) GetTokenNo() lex.Token { return tn.Token }
func (tn *TipoNone) GetInformacao() string { return ferramentas.GetIdentacao(tn.Profu) + "Tipo NONE" }
