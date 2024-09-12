package arvore

import (
	"ChikenInterpreter/ferramentas"
	lex "ChikenInterpreter/lexing"
	"fmt"
	"strings"
)

// definicoes de interfaces
type No interface {
	GetTokenNo() lex.Token
	GetInformacao() string
}

type Instrucao interface {
	No
	noInstrucao()
	IncProfu()
}

type Expressao interface {
	No
	noExpressao()
	IncProfu()
}

// Blocos de instrucao
type BlocoInstrucao struct {
	Token      lex.Token
	Instrucoes []Instrucao
}

func (bi *BlocoInstrucao) noInstrucao() {}
func (bi *BlocoInstrucao) IncProfu() {
	for _, i := range bi.Instrucoes {
		i.IncProfu()
	}
}
func (bi *BlocoInstrucao) GetTokenNo() lex.Token {
	return bi.Token
}
func (bi *BlocoInstrucao) GetInformacao() string {
	if len(bi.Instrucoes) > 0 {
		var parts []string = make([]string, len(bi.Instrucoes))

		for i, instru := range bi.Instrucoes {
			parts[i] = instru.GetInformacao()
		}

		return strings.Join(parts, "\n")
	}

	return "Bloco vazio"
}

// estrutura do programa
type Programa struct {
	Instrucoes []Instrucao
}

func (p *Programa) GetTokenNo() lex.Token {
	if len(p.Instrucoes) > 0 {
		return p.Instrucoes[0].GetTokenNo()
	}

	return lex.Token{Tipo: lex.ILEGAL, Valor: ""}
}

func (p *Programa) GetInformacao() (texto string) {
	texto = "Programa vazio"

	if len(p.Instrucoes) > 0 {
		var parts []string

		for _, instru := range p.Instrucoes {
			parts = append(parts, "\n"+instru.GetInformacao())
		}
		texto = strings.Join(parts, "")
		return
	}

	return
}

// Tipos de Instrucao:
// instrucao de expressao
type InstrucaodeExpressao struct {
	Token     lex.Token
	Expressao Expressao
	Profu     int
}

func (ie *InstrucaodeExpressao) noInstrucao() {}
func (ie *InstrucaodeExpressao) IncProfu() {
	ie.Profu++
	ie.Expressao.IncProfu()
}
func (ie *InstrucaodeExpressao) GetTokenNo() lex.Token {
	return ie.Token
}
func (ie *InstrucaodeExpressao) GetInformacao() string {
	return fmt.Sprintf("%sInstrucao de expressao:\n%s", ferramentas.GetIdentacao(ie.Profu), ie.Expressao.GetInformacao())
}

// instrucao de declaracao de variavel
type VarInstrucao struct {
	Token lex.Token
	Vars  []*VarDeclaracao
	Profu int
}

func (vi *VarInstrucao) noInstrucao() {}
func (vi *VarInstrucao) IncProfu() {
	for _, v := range vi.Vars {
		v.IncProfu()
	}
}
func (vi *VarInstrucao) GetTokenNo() lex.Token { return vi.Token }
func (vi *VarInstrucao) GetInformacao() string {
	indent := ferramentas.GetIdentacao(vi.Profu)
	declaracoes := make([]string, len(vi.Vars)+1)

	declaracoes[0] = indent + "Declaracao de variaveis:\n"
	for i, variavel := range vi.Vars {
		declaracoes[i+1] = variavel.GetInformacao()
	}

	return strings.Join(declaracoes, "\n")
}

type VarDeclaracao struct {
	Token  lex.Token
	Ident  *Identificador
	Expres Expressao
	Profu  int
}

func (v *VarDeclaracao) noInstrucao() {}
func (v *VarDeclaracao) IncProfu() {
	v.Expres.IncProfu()
	v.Profu++
}
func (v *VarDeclaracao) GetTokenNo() lex.Token {
	return v.Token
}
func (v *VarDeclaracao) GetInformacao() string {
	indent := ferramentas.GetIdentacao(v.Profu)
	return fmt.Sprintf(indent+"Token: %s, %s, Expressao:\n%s", string(v.Token.Tipo), v.Ident.Token.Valor, v.Expres.GetInformacao())
}

// Instrucao def
type DefInstrucao struct {
	Token  lex.Token
	Ident  *Identificador
	Expres Expressao
	Profu  int
}

func (di *DefInstrucao) noInstrucao() {}
func (di *DefInstrucao) IncProfu() {
	di.Expres.IncProfu()
	di.Profu++
}
func (di *DefInstrucao) GetTokenNo() lex.Token {
	return di.Token
}
func (di *DefInstrucao) GetInformacao() string {
	indent := ferramentas.GetIdentacao(di.Profu)
	return fmt.Sprintf(indent+"Token: %s, %s, Expressao:\n%s", string(di.Token.Tipo), di.Ident.Token.Valor, di.Expres.GetInformacao())
}

// instrucao de return
type ReturnInstrucao struct {
	Token lex.Token
	Expre Expressao
	Profu int
}

func (r *ReturnInstrucao) noInstrucao() {}
func (r *ReturnInstrucao) IncProfu() {
	r.Expre.IncProfu()
	r.Profu++
}
func (r *ReturnInstrucao) GetTokenNo() lex.Token {
	return r.Token
}
func (r *ReturnInstrucao) GetInformacao() string {
	indent := ferramentas.GetIdentacao(r.Profu)
	texto := indent + "Instrucao RETURN, expressao:\n"
	if r.Expre == nil {
		texto += "Sem expressao"
	} else {
		texto += r.Expre.GetInformacao()
	}

	return texto
}

// Instrucao de run
type RunInstrucao struct {
	Token lex.Token
	Expre Expressao
	Profu int
}

func (ri *RunInstrucao) noInstrucao() {}
func (ri *RunInstrucao) IncProfu() {
	ri.Expre.IncProfu()
	ri.Profu++
}
func (ri *RunInstrucao) GetTokenNo() lex.Token {
	return ri.Token
}
func (ri *RunInstrucao) GetInformacao() string {
	indent := ferramentas.GetIdentacao(ri.Profu)
	texto := indent + "Instrucao RUN, expressao:\n"
	texto += ri.Expre.GetInformacao()

	return texto
}

// Instrucao super
type SuperInstrucao struct {
	Token lex.Token
	ClasseMae string
	Propriedade string
	Argumentos []Expressao
	Profu int
}

func (si *SuperInstrucao) noInstrucao() {}
func (si *SuperInstrucao) IncProfu() {
	for _,i := range si.Argumentos {
		i.IncProfu()
	}
	si.Profu++
}
func (si *SuperInstrucao) GetTokenNo() lex.Token { return si.Token }
func (si *SuperInstrucao) GetInformacao() string {
	indent := ferramentas.GetIdentacao(si.Profu)
	texto := indent + fmt.Sprintf("Instrucao SUPER. Classe chamada: %s. Atributo recuperado: %s\n",si.ClasseMae,si.Propriedade)

	if len(si.Argumentos) > 0{
		texto += indent+"Argumentos:\n"

		for _,i := range si.Argumentos{
			texto += i.GetInformacao()+"\n"
		}
	}

	return texto
}

// Instrucao de erro
type ErrInstrucao struct {
	Token lex.Token
	Expre Expressao
	Profu int
}

func (ei *ErrInstrucao) noInstrucao() {}
func (ei *ErrInstrucao) IncProfu() {
	ei.Expre.IncProfu()
	ei.Profu++
}
func (ei *ErrInstrucao) GetTokenNo() lex.Token {
	return ei.Token
}
func (ei *ErrInstrucao) GetInformacao() string {
	indent := ferramentas.GetIdentacao(ei.Profu)
	texto := indent + "Instrucao ERROR, expressao:\n"
	texto += ei.Expre.GetInformacao()

	return texto
}

type InstrucaoBreak struct {
	Token lex.Token
	Profu int
}

func (ib *InstrucaoBreak) noInstrucao()          {}
func (ib *InstrucaoBreak) IncProfu()             { ib.Profu++ }
func (ib *InstrucaoBreak) GetTokenNo() lex.Token { return ib.Token }
func (ib *InstrucaoBreak) GetInformacao() string {
	return ferramentas.GetIdentacao(ib.Profu) + "Instrução BREAK"
}

type InstrucaoContinue struct {
	Token lex.Token
	Profu int
}

func (ic *InstrucaoContinue) noInstrucao()          {}
func (ic *InstrucaoContinue) IncProfu()             { ic.Profu++ }
func (ic *InstrucaoContinue) GetTokenNo() lex.Token { return ic.Token }
func (ic *InstrucaoContinue) GetInformacao() string {
	return ferramentas.GetIdentacao(ic.Profu) + "Instrução CONTINUE"
}

type InstrucaoAtribuicao struct {
	Token      lex.Token
	Operador   string
	ExprRecebe Expressao
	ExprValue  Expressao
	Profu      int
}

func (ia *InstrucaoAtribuicao) noInstrucao() {}
func (ia *InstrucaoAtribuicao) IncProfu() {
	ia.ExprRecebe.IncProfu()
	ia.ExprValue.IncProfu()
}
func (ia *InstrucaoAtribuicao) GetTokenNo() lex.Token { return ia.Token }
func (ia *InstrucaoAtribuicao) GetInformacao() string {
	parts := make([]string, 3)
	indent := ferramentas.GetIdentacao(ia.Profu)

	parts[0] = indent + "Atribuicao de valor, operador: " + ia.Operador
	parts[1] = indent + "Receptor:\n" + ia.ExprRecebe.GetInformacao()
	parts[2] = indent + "Valor:\n" + ia.ExprValue.GetInformacao()

	return strings.Join(parts, "\n")
}

// Instrucoes complexas
type InstrucaoIter struct {
	Token          lex.Token
	Iterador       *Identificador
	ExpressaoLista Expressao
	BlocoCodigo    *BlocoInstrucao
	Profu          int
}

func (ii *InstrucaoIter) noInstrucao() {}
func (ii *InstrucaoIter) IncProfu() {
	ii.Iterador.IncProfu()
	ii.ExpressaoLista.IncProfu()
	ii.BlocoCodigo.IncProfu()
	ii.Profu++
}
func (ii *InstrucaoIter) GetTokenNo() lex.Token {
	return ii.Token
}
func (ii *InstrucaoIter) GetInformacao() string {
	parts := make([]string, 7)
	indent := ferramentas.GetIdentacao(ii.Profu)

	parts[0] = indent + "Instrucao de iteracao."
	parts[1] = indent + "Iterador:"
	parts[2] = ii.Iterador.GetInformacao()
	parts[3] = indent + "Expressao de lista:"
	parts[4] = ii.ExpressaoLista.GetInformacao()
	parts[5] = indent + "Codigo:"
	parts[6] = ii.BlocoCodigo.GetInformacao()

	return strings.Join(parts, "\n")
}

type InstrucaoSwitch struct {
	Token        lex.Token
	ExpreTeste   Expressao
	Cases        []*Case
	BlocoDefault *BlocoInstrucao
	Profu        int
}

func (is *InstrucaoSwitch) noInstrucao() {}
func (is *InstrucaoSwitch) IncProfu() {
	is.ExpreTeste.IncProfu()
	for _, caso := range is.Cases {
		caso.ExpreCase.IncProfu()
	}

	is.Profu++
}

func (is *InstrucaoSwitch) GetTokenNo() lex.Token {
	return is.Token
}
func (is *InstrucaoSwitch) GetInformacao() string {
	parts := make([]string, len(is.Cases)+3)
	indent := ferramentas.GetIdentacao(is.Profu)

	parts[0] = indent + "Instrucao SWITCH.\n" + indent + "Expressao de entrada:"
	parts[1] = is.ExpreTeste.GetInformacao()

	for i, caso := range is.Cases {
		parts[i+2] = indent + caso.GetInformacao() + "\n" + indent + "Codigo:\n" + caso.Codigo.GetInformacao()
	}

	if is.BlocoDefault == nil {
		return strings.Join(parts,"\n")
	}
	parts[len(parts)-1] = indent + "Default:\n" + is.BlocoDefault.GetInformacao()

	return strings.Join(parts, "\n")
}

type Case struct {
	ExpreCase Expressao
	Codigo    *BlocoInstrucao
}

func (c *Case) noExpressao() {}
func (c *Case) IncProfu() {
	c.ExpreCase.IncProfu()
	c.Codigo.IncProfu()
}
func (c *Case) GetTokenNo() lex.Token {
	return lex.Token{Valor: "case", Tipo: "CASE"}
}
func (c *Case) GetInformacao() string {
	return "Caso:\n" + c.ExpreCase.GetInformacao()
}

type InstrucaoTryExcept struct {
	Token       lex.Token
	ExcessaoVar string
	BlocoTry    *BlocoInstrucao
	BlocoExcept *BlocoInstrucao
	Profu       int
}

func (ite *InstrucaoTryExcept) noInstrucao() {}
func (ite *InstrucaoTryExcept) IncProfu() {
	ite.BlocoTry.IncProfu()
	ite.BlocoExcept.IncProfu()
	ite.Profu++
}
func (ite *InstrucaoTryExcept) GetTokenNo() lex.Token {
	return ite.Token
}
func (ite *InstrucaoTryExcept) GetInformacao() string {
	indent := ferramentas.GetIdentacao(ite.Profu)
	parts := make([]string, 6)

	parts[0] = indent + "Instrução try-except."
	parts[1] = indent + "Variavel e captura de exceção: " + ite.ExcessaoVar

	parts[2] = indent + "Bloco try:"
	parts[3] = ite.BlocoTry.GetInformacao()

	parts[4] = indent + "Bloco except:"
	parts[5] = ite.BlocoExcept.GetInformacao()

	return strings.Join(parts, "\n")
}
