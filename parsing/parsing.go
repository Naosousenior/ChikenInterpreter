package parsing

import (
	lex "ChikenInterpreter/lexing"
	arv "ChikenInterpreter/parsing/arvore"
	aux "ChikenInterpreter/parsing/auxiliar"
	"fmt"
	"strconv"
)

type (
	operadorPrefixo func() (arv.Expressao, bool)
	operadorInfixo  func(arv.Expressao) (arv.Expressao, bool)
	//operadorPosfixo func() Expressao
)

const (
	_          int = iota
	ATRIBUICAO     //atribuicao
	MENOR          //definicao
	LOG_E          //And
	LOG_OU         //Or
	IN             //IN
	COMP           //==
	MA_ME          //< e >
	SOMA_SUB       //+ e -
	MUL_DIV        //* e / e %
	POW            //**
	PREFIX         //! e -
	FUNCAO         //funcao(parametros)
	ATRIBUTO       // .
)

var precedencias = map[lex.TipoToken]int{
	lex.COMPARACAO_IG:   COMP,
	lex.COMPARACAO_DIF:  COMP,
	lex.RECEBE:          ATRIBUICAO,
	lex.ADD_RECEBE:      ATRIBUICAO,
	lex.SUB_RECEBE:      ATRIBUICAO,
	lex.MUL_RECEBE:      ATRIBUICAO,
	lex.DIV_RECEBE:      ATRIBUICAO,
	lex.MOD_RECEBE:      ATRIBUICAO,
	lex.TIPO_RECEBE:     ATRIBUICAO,
	lex.DOIS_PONTO:      IN,
	lex.TESTE_IS:        COMP,
	lex.MENOR_Q:         MA_ME,
	lex.MENOR_IGUAL:     MA_ME,
	lex.MAIOR_Q:         MA_ME,
	lex.MAIOR_IGUAL:     MA_ME,
	lex.AND:             LOG_E,
	lex.OR:              LOG_OU,
	lex.XOR:             LOG_OU,
	lex.ADD:             SOMA_SUB,
	lex.SUB:             SOMA_SUB,
	lex.MUL:             MUL_DIV,
	lex.DIV:             MUL_DIV,
	lex.RESTO:           MUL_DIV,
	lex.POTENCIA:        POW,
	lex.PONTO:           ATRIBUTO,
	lex.ABRE_PARENTESES: FUNCAO,
	lex.ABRE_COLCHETE:   ATRIBUTO,
}

type Analisador struct {
	lexico *lex.Lexico

	erros []string
	profu int

	Atual   lex.Token
	Proximo lex.Token

	opsPrefixoFn map[lex.TipoToken]operadorPrefixo
	opsInfixoFn  map[lex.TipoToken]operadorInfixo
}

func NewAnalisador(lexico *lex.Lexico) *Analisador {
	a := &Analisador{lexico: lexico, erros: []string{}, profu: 0}
	a.opsPrefixoFn = make(map[lex.TipoToken]operadorPrefixo)
	a.opsInfixoFn = make(map[lex.TipoToken]operadorInfixo)
	a.avancaToken()
	a.avancaToken()

	a.addPrefixos()
	a.addInfixos()

	return a
}

func (a *Analisador) addPrefixos() {
	a.addPrefixoFn(lex.IDENTIFICADOR, a.parseIdentificador)
	a.addPrefixoFn(lex.OBJECT, a.parseExpressaoObjeto)
	a.addPrefixoFn(lex.CLASS, a.parseExpressaoClass)
	a.addPrefixoFn(lex.NUM_INT, a.parseLiteralInteiro)
	a.addPrefixoFn(lex.NUM_REAL, a.parseLiteralReal)
	a.addPrefixoFn(lex.TRUE, a.parseBooleano)
	a.addPrefixoFn(lex.FALSE, a.parseBooleano)
	a.addPrefixoFn(lex.STRING, a.parseString)
	a.addPrefixoFn(lex.NONE, a.parseTipoNone)
	a.addPrefixoFn(lex.NEGACAO, a.parseExpressaoPrefixo)
	a.addPrefixoFn(lex.SUB, a.parseExpressaoPrefixo)
	a.addPrefixoFn(lex.ABRE_COLCHETE, a.parseExpressaoLista)
	a.addPrefixoFn(lex.ABRE_PARENTESES, a.parseExpressaoAgrupada)
	a.addPrefixoFn(lex.ABRE_CHAVE, a.parseExpressaoDict)
	a.addPrefixoFn(lex.IF, a.parseExpressaoIf)
	a.addPrefixoFn(lex.REPEAT, a.parseExpressaoRepeat)
	a.addPrefixoFn(lex.FUN, a.parseFunExpressao)
}

func (a *Analisador) addInfixos() {
	a.addInfixoFn(lex.ADD, a.parseExpressaoInfixo)
	a.addInfixoFn(lex.SUB, a.parseExpressaoInfixo)
	a.addInfixoFn(lex.MUL, a.parseExpressaoInfixo)
	a.addInfixoFn(lex.DIV, a.parseExpressaoInfixo)
	a.addInfixoFn(lex.POTENCIA, a.parseExpressaoInfixo)
	a.addInfixoFn(lex.RESTO, a.parseExpressaoInfixo)
	a.addInfixoFn(lex.MENOR_Q, a.parseExpressaoInfixo)
	a.addInfixoFn(lex.MENOR_IGUAL, a.parseExpressaoInfixo)
	a.addInfixoFn(lex.MAIOR_Q, a.parseExpressaoInfixo)
	a.addInfixoFn(lex.MAIOR_IGUAL, a.parseExpressaoInfixo)
	a.addInfixoFn(lex.COMPARACAO_IG, a.parseExpressaoInfixo)
	a.addInfixoFn(lex.COMPARACAO_DIF, a.parseExpressaoInfixo)
	a.addInfixoFn(lex.DOIS_PONTO, a.parseExpressaoInfixo)
	a.addInfixoFn(lex.TESTE_IS, a.parseExpressaoInfixo)
	a.addInfixoFn(lex.AND, a.parseExpressaoInfixo)
	a.addInfixoFn(lex.OR, a.parseExpressaoInfixo)
	a.addInfixoFn(lex.XOR, a.parseExpressaoInfixo)
	a.addInfixoFn(lex.PONTO, a.parseExpressaoAtributo)
	a.addInfixoFn(lex.ABRE_PARENTESES, a.parseChamadaFuncao)
	a.addInfixoFn(lex.ABRE_COLCHETE, a.parseExpressaoIndex)
}

func (a *Analisador) GetErros() []string {
	return a.erros
}

func (a *Analisador) addErro(erro string) {
	posicao := a.Atual.StringPos()
	texto_erro := fmt.Sprintf("Erro em: %s. %s", posicao, erro)
	a.erros = append(a.erros, texto_erro)
}

func (a *Analisador) addPrefixoFn(tipo lex.TipoToken, funcao operadorPrefixo) {
	a.opsPrefixoFn[tipo] = funcao
}

func (a *Analisador) addInfixoFn(tipo lex.TipoToken, funcao operadorInfixo) {
	a.opsInfixoFn[tipo] = funcao
}

func (a *Analisador) avancaToken() {
	a.Atual = a.Proximo
	a.Proximo = a.lexico.GetToken()
}

func (a *Analisador) analisaInstrucao() arv.Instrucao {
	switch a.Atual.Tipo {
	case lex.VAR:
		return a.parseVar()

	case lex.DEFINITION:
		return a.parseDef()

	case lex.IF:
		return a.parseIf()

	case lex.SWITCH:
		return a.parseSwitch()

	case lex.REPEAT:
		return a.parseRepeat()

	case lex.BREAK:
		instru := &arv.InstrucaoBreak{Token: a.Atual, Profu: a.profu}
		if a.esperaToken(lex.PONTO_VIRGULA) {
			return instru
		} else {
			return nil
		}

	case lex.CONTINUE:
		instru := &arv.InstrucaoContinue{Token: a.Atual, Profu: a.profu}
		if a.esperaToken(lex.PONTO_VIRGULA) {
			return instru
		} else {
			return nil
		}

	case lex.ITER:
		return a.parseIter()

	case lex.RETURN:
		return a.parseReturn()

	case lex.RUN:
		return a.parseRun()

	case lex.SUPER:
		return a.parseSuper()

	case lex.ERROR:
		return a.parseErr()
	case lex.TRY:
		return a.parseTryExcept()

	default:
		return a.parseInstrucaoExpressao()
	}
}

func (a *Analisador) AnalisaPrograma() *arv.Programa {
	programa := &arv.Programa{}
	programa.Instrucoes = []arv.Instrucao{}

	for a.Atual.Tipo != lex.FIM {

		instrucao := a.analisaInstrucao()
		if instrucao != nil {
			programa.Instrucoes = append(programa.Instrucoes, instrucao)
		}

		a.avancaToken()
	}

	return programa
}

// parses:
func (a *Analisador) esperaToken(t lex.TipoToken) bool {
	if a.Proximo.Tipo == t {
		a.avancaToken()
		return true
	}

	a.addErro(fmt.Sprintf("Esperavamos o Token %s, mas tivemos %s", string(t), string(a.Proximo.Tipo)))
	return false
}

func (a *Analisador) testaTokenAtual(t lex.TipoToken) bool {
	return a.Atual.Tipo == t
}

func (a *Analisador) testaProximoToken(t lex.TipoToken) bool {
	return a.Proximo.Tipo == t
}

func (a *Analisador) prioridadeProximo() int {
	return precedencias[a.Proximo.Tipo]
}

func (a *Analisador) prioridadeAtual() int {
	return precedencias[a.Atual.Tipo]
}

func (a *Analisador) parseVar() arv.Instrucao {
	instru := &arv.VarInstrucao{Token: a.Atual, Profu: a.profu}
	instru.Vars = []*arv.VarDeclaracao{}
	a.profu++
	for a.testaProximoToken(lex.IDENTIFICADOR) {
		a.avancaToken()
		dec := a.parseDeclaracao()

		if dec == nil {
			return nil
		}

		if a.testaProximoToken(lex.VIRGULA) {
			a.avancaToken()
		}

		instru.Vars = append(instru.Vars, dec)
	}

	a.profu--

	if len(instru.Vars) > 0 {
		if !a.esperaToken(lex.PONTO_VIRGULA) {
			return nil
		}

		return instru
	}

	a.addErro("Nenhuma variavel corretamente declarada")

	return nil
}

func (a *Analisador) parseDef() arv.Instrucao {
	instru := &arv.DefInstrucao{Token: a.Atual, Profu: a.profu}

	if a.esperaToken(lex.IDENTIFICADOR) {
		var ok bool
		instru.Ident = &arv.Identificador{Profu: a.profu + 1, Token: a.Atual, Nome: a.Atual.Valor}

		a.avancaToken()
		instru.Expres, ok = a.parseExpressao(MENOR)

		if !(ok && a.esperaToken(lex.PONTO_VIRGULA)) {
			return nil
		}

		return instru
	}

	return nil
}

func (a *Analisador) parseDeclaracao() *arv.VarDeclaracao {
	declaracao := arv.VarDeclaracao{Token: a.Atual, Profu: a.profu}
	var err bool

	declaracao.Ident = &arv.Identificador{Token: a.Atual, Nome: a.Atual.Valor, Profu: a.profu}

	if a.testaProximoToken(lex.RECEBE) {
		a.avancaToken()
		a.avancaToken()

		declaracao.Expres, err = a.parseExpressao(MENOR)

		if !err {
			return nil
		}

	} else if a.testaProximoToken(lex.PONTO_VIRGULA) || a.testaProximoToken(lex.VIRGULA) {

		declaracao.Expres = &arv.TipoNone{Token: lex.Token{Tipo: lex.NONE, Valor: "none"}, Profu: a.profu}
	} else {
		a.addErro(fmt.Sprintf("Esperavamos \"%s\",\"%s\", ou \"%s\", mas tivemos: %s", lex.RECEBE, lex.VIRGULA, lex.PONTO_VIRGULA, a.Proximo.Tipo))
		return nil
	}

	return &declaracao
}

func (a *Analisador) parseReturn() arv.Instrucao {
	instrucao := arv.ReturnInstrucao{Token: a.Atual, Profu: a.profu}
	var err bool
	a.avancaToken()
	if a.testaTokenAtual(lex.PONTO_VIRGULA) {
		return &instrucao
	}

	instrucao.Expre, err = a.parseExpressao(MENOR)

	if !err {
		return nil
	}

	a.esperaToken(lex.PONTO_VIRGULA)

	return &instrucao
}

func (a *Analisador) parseRun() arv.Instrucao {
	instru := arv.RunInstrucao{Token: a.Atual, Profu: a.profu}
	var err bool
	a.avancaToken()
	instru.Expre, err = a.parseExpressao(MENOR)

	if !err {
		return nil
	}

	a.esperaToken(lex.PONTO_VIRGULA)

	return &instru
}

func (a *Analisador) parseSuper() arv.Instrucao {
	instru := arv.SuperInstrucao{Token: a.Atual, Profu: a.profu}
	instru.Argumentos = make([]arv.Expressao, 0)

	if !a.esperaToken(lex.IDENTIFICADOR) {
		return nil
	}

	instru.ClasseMae = a.Atual.Valor

	if !a.esperaToken(lex.PONTO) {
		return nil
	}

	if !a.esperaToken(lex.IDENTIFICADOR) {
		return nil
	}

	instru.Propriedade = a.Atual.Valor

	if a.testaProximoToken(lex.ABRE_PARENTESES) {
		a.avancaToken()

		instru.Argumentos = a.parseArgumentosChamada()
	}

	return &instru
}

func (a *Analisador) parseErr() arv.Instrucao {
	instru := arv.ErrInstrucao{Token: a.Atual, Profu: a.profu}
	var err bool

	a.avancaToken()
	instru.Expre, err = a.parseExpressao(MENOR)

	if !err {
		return nil
	}

	a.esperaToken(lex.PONTO_VIRGULA)

	return &instru
}

func (a *Analisador) parseTryExcept() arv.Instrucao {
	instru := &arv.InstrucaoTryExcept{Token: a.Atual, Profu: a.profu}
	var blocoTry, blocoErr *arv.BlocoInstrucao

	if !a.esperaToken(lex.ABRE_CHAVE) {
		return nil
	}

	blocoTry = a.parseBlocoInstrucoes()

	if !a.esperaToken(lex.EXCEPT) {
		return nil
	}

	if a.testaProximoToken(lex.IDENTIFICADOR) {
		a.avancaToken()

		instru.ExcessaoVar = a.Atual.Valor
	}

	if !a.esperaToken(lex.ABRE_CHAVE) {
		return nil
	}

	blocoErr = a.parseBlocoInstrucoes()

	instru.BlocoTry = blocoTry
	instru.BlocoExcept = blocoErr

	return instru
}

func (a *Analisador) parseInstrucaoExpressao() arv.Instrucao {
	instrucao := &arv.InstrucaodeExpressao{Token: a.Atual, Profu: a.profu}
	var err bool

	instrucao.Expressao, err = a.parseExpressao(MENOR)
	if !err {
		return nil
	} else if a.Proximo.IsTokenRecebe() {
		a.avancaToken()
		var expValor arv.Expressao
		instru := &arv.InstrucaoAtribuicao{Token: a.Atual,
			Operador:   a.Atual.Valor,
			ExprRecebe: instrucao.Expressao,
			Profu:      a.profu}

		a.avancaToken()
		expValor, err = a.parseExpressao(MENOR)

		if !err {
			return nil
		}

		instru.ExprValue = expValor

		a.esperaToken(lex.PONTO_VIRGULA)

		return instru
	}

	a.esperaToken(lex.PONTO_VIRGULA)

	return instrucao
}

func (a *Analisador) parseIf() arv.Instrucao {
	ifInstrucao := &arv.InstrucaodeExpressao{Token: a.Atual, Profu: a.profu}
	ifExpressao, ok := a.parseExpressaoIf()

	if ok {
		ifInstrucao.Expressao = ifExpressao
		return ifInstrucao
	}

	return nil
}

func (a *Analisador) parseSwitch() arv.Instrucao {
	switInstru := &arv.InstrucaoSwitch{Token: a.Atual, Profu: a.profu}
	filaCases := aux.NewFilaCasos()
	var ok bool

	if !a.esperaToken(lex.ABRE_PARENTESES) {
		return nil
	}

	a.avancaToken()
	switInstru.ExpreTeste, ok = a.parseExpressao(MENOR)
	if !(ok && a.esperaToken(lex.FECHA_PARENTESES)) {
		return nil
	}

	if !a.esperaToken(lex.ABRE_CHAVE) {
		return nil
	}

	a.profu++
	for a.testaProximoToken(lex.CASE) {
		a.avancaToken()
		a.profu++
		novoCaso := a.parseCaso()
		a.profu--

		if novoCaso == nil {
			a.profu--
			return nil
		}
		filaCases.AddItem(novoCaso)
	}
	a.profu++

	if a.testaProximoToken(lex.DEFAULT) {
		a.avancaToken()
		a.avancaToken()
		switInstru.BlocoDefault = a.parseBlocoInstrucoes()
	}

	if !a.esperaToken(lex.FECHA_CHAVE) {
		return nil
	}

	switInstru.Cases = make([]*arv.Case, filaCases.Tamanho)
	for i := 0; i < filaCases.Tamanho; i++ {
		switInstru.Cases[i] = filaCases.GetItem()
	}

	return switInstru
}

func (a *Analisador) parseCaso() *arv.Case {
	exprCaso := &arv.Case{}
	var ok bool

	a.avancaToken()

	exprCaso.ExpreCase, ok = a.parseExpressao(MENOR)

	if !ok {
		return nil
	}

	if !a.esperaToken(lex.ABRE_CHAVE) {
		return nil
	}

	exprCaso.Codigo = a.parseBlocoInstrucoes()

	return exprCaso
}

func (a *Analisador) parseRepeat() arv.Instrucao {
	rptInstru := &arv.InstrucaodeExpressao{Token: a.Atual, Profu: a.profu}
	repeatExpressao, ok := a.parseExpressaoRepeat()

	if ok {
		rptInstru.Expressao = repeatExpressao
		return rptInstru
	}

	return nil
}

func (a *Analisador) parseIter() arv.Instrucao {
	instru := &arv.InstrucaoIter{Token: a.Atual, Profu: a.profu}
	var ok bool

	if !a.esperaToken(lex.ABRE_PARENTESES) {
		return nil
	}

	if !a.esperaToken(lex.IDENTIFICADOR) {
		return nil
	}

	instru.Iterador = &arv.Identificador{Token: a.Atual, Nome: a.Atual.Valor, Profu: a.profu + 1}

	if !a.esperaToken(lex.SETA) {
		return nil
	}

	a.avancaToken()

	instru.ExpressaoLista, ok = a.parseExpressao(MENOR)

	if !(ok && a.esperaToken(lex.FECHA_PARENTESES)) {
		return nil
	}

	if !a.esperaToken(lex.ABRE_CHAVE) {
		return nil
	}

	instru.BlocoCodigo = a.parseBlocoInstrucoes()

	return instru
}

func (a *Analisador) parseExpressao(prioridade int) (arv.Expressao, bool) {
	prefixFun := a.opsPrefixoFn[a.Atual.Tipo]

	if prefixFun == nil {
		a.addErro(fmt.Sprintf("Token %s inesperado ou desconhecido", a.Atual.Valor))
		return nil, false
	}

	a.profu++
	expressao, err := prefixFun()
	if !err {
		return nil, false
	}

	for !a.testaProximoToken(lex.PONTO_VIRGULA) && prioridade < a.prioridadeProximo() {
		infixoFun := a.opsInfixoFn[a.Proximo.Tipo]
		if infixoFun == nil {
			a.addErro(fmt.Sprintf("Operador %s nao reconhecido", a.Proximo.Valor))
			return nil, false
		}

		a.avancaToken()

		expressao.IncProfu()
		expressao, err = infixoFun(expressao)
		if !err {
			return nil, false
		}
	}

	a.profu--
	return expressao, true
}

func (a *Analisador) parseExpressaoIf() (arv.Expressao, bool) {
	var err bool
	expr := &arv.ExpressaoIf{Token: a.Atual, Profu: a.profu}

	if !a.esperaToken(lex.ABRE_PARENTESES) {
		return nil, false
	}

	a.avancaToken()
	expr.Condicao, err = a.parseExpressao(MENOR)

	if !err {
		return nil, false
	}

	if !a.esperaToken(lex.FECHA_PARENTESES) {
		return nil, false
	}
	if !a.esperaToken(lex.ABRE_CHAVE) {
		return nil, false
	}

	expr.BlocoEntao = a.parseBlocoInstrucoes()

	if a.testaProximoToken(lex.ELSE) {
		a.avancaToken()
		if a.testaProximoToken(lex.IF) {
			instruExpr := &arv.InstrucaodeExpressao{Token: a.Proximo, Profu: a.profu}
			a.avancaToken()
			elseif, err := a.parseExpressao(POW)

			if err {
				instruExpr.Expressao = elseif
				expr.BlocoSenao = &arv.BlocoInstrucao{Instrucoes: []arv.Instrucao{}}
				expr.BlocoSenao.Instrucoes = append(expr.BlocoSenao.Instrucoes, instruExpr)

				return expr, true
			} else {
				return nil, false
			}
		} else if !a.esperaToken(lex.ABRE_CHAVE) {
			return nil, false
		}
		expr.BlocoSenao = a.parseBlocoInstrucoes()
	}

	return expr, true
}

func (a *Analisador) parseExpressaoRepeat() (arv.Expressao, bool) {
	rptExpr := &arv.ExpressaoRepeat{Token: a.Atual, Profu: a.profu}
	var ok bool

	if a.testaProximoToken(lex.ABRE_PARENTESES) {
		a.avancaToken()
		a.avancaToken()
		rptExpr.Condicao1, ok = a.parseExpressao(MENOR)
		if !(ok && a.esperaToken(lex.FECHA_PARENTESES)) {
			return nil, false
		}

	} else {
		rptExpr.Condicao1 = &arv.Booleano{Token: lex.Token{Valor: "true", Tipo: lex.TRUE}, Valor: true, Profu: a.profu + 1}
	}

	if !a.esperaToken(lex.ABRE_CHAVE) {
		return nil, false
	}

	rptExpr.BlocoRepetir = a.parseBlocoInstrucoes()

	if a.testaProximoToken(lex.ABRE_PARENTESES) {
		a.avancaToken()
		a.avancaToken()
		rptExpr.Condicao2, ok = a.parseExpressao(MENOR)
		if !(ok && a.esperaToken(lex.FECHA_PARENTESES)) {
			return nil, false
		}
	} else {
		rptExpr.Condicao2 = &arv.Booleano{Token: lex.Token{Valor: "true", Tipo: lex.TRUE}, Valor: true, Profu: a.profu + 1}
	}

	return rptExpr, true
}

func (a *Analisador) parseFunExpressao() (arv.Expressao, bool) {
	exprFun := &arv.ExpressaoFun{Token: a.Atual, Profu: a.profu}

	if !a.esperaToken(lex.ABRE_PARENTESES) {
		return nil, false
	}

	a.avancaToken()
	exprFun.Parametros = a.parseParametros()

	if exprFun.Parametros == nil {
		return nil, false
	}

	if !a.esperaToken(lex.ABRE_CHAVE) {
		return nil, false
	}

	exprFun.Bloco = a.parseBlocoInstrucoes()

	return exprFun, true
}

func (a *Analisador) parseChamadaFuncao(expEsque arv.Expressao) (arv.Expressao, bool) {
	expre := &arv.CallFun{Funcao: expEsque, Profu: a.profu}
	expre.Argumentos = a.parseArgumentosChamada()

	if expre.Argumentos == nil {
		return nil, false
	}

	return expre, true
}

func (a *Analisador) parseParametros() []*arv.Identificador {
	identi := []*arv.Identificador{}

	for a.testaTokenAtual(lex.IDENTIFICADOR) {
		identi = append(identi, &arv.Identificador{Token: a.Atual, Nome: a.Atual.Valor, Profu: a.profu + 1})

		if a.testaProximoToken(lex.VIRGULA) {
			a.avancaToken()
			a.avancaToken()
		} else {
			a.avancaToken()
			break
		}
	}

	if !a.testaTokenAtual(lex.FECHA_PARENTESES) {
		a.addErro(fmt.Sprintf("Os parametros da funcao devem ser limitados por parenteses, nao por %s", a.Atual.Tipo))
		return nil
	}

	return identi
}

func (a *Analisador) parseArgumentosChamada() []arv.Expressao {
	argumentos := []arv.Expressao{}

	if a.testaProximoToken(lex.FECHA_PARENTESES) {
		a.avancaToken()
		return argumentos
	}

	a.avancaToken()
	arg, err := a.parseExpressao(MENOR)

	if !err {
		return nil
	}

	filaArgs := aux.NewFilaExpressoes()

	filaArgs.AddItem(arg)

	for a.testaProximoToken(lex.VIRGULA) {
		a.avancaToken()
		a.avancaToken()
		arg, err = a.parseExpressao(MENOR)
		if !err {
			return nil
		}

		filaArgs.AddItem(arg)
	}

	if !a.esperaToken(lex.FECHA_PARENTESES) {
		return nil
	}

	argumentos = aux.GetLista(filaArgs)

	return argumentos
}

func (a *Analisador) parseExpressaoIndex(expr arv.Expressao) (arv.Expressao, bool) {
	expressaoIndex := &arv.ExpressaoInfixo{Token: a.Atual, Profu: a.profu, Operador: a.Atual.Valor, ExpEsquerda: expr}
	a.avancaToken()
	index, ok := a.parseExpressao(MENOR)
	if !ok {
		return nil, false
	}

	expressaoIndex.ExpDireita = index

	if !a.esperaToken(lex.FECHA_COLCHETE) {
		return nil, false
	}

	return expressaoIndex, true
}

func (a *Analisador) parseBlocoInstrucoes() *arv.BlocoInstrucao {
	bloco := &arv.BlocoInstrucao{Token: a.Atual}
	bloco.Instrucoes = []arv.Instrucao{}

	a.avancaToken()

	a.profu++
	for !(a.testaTokenAtual(lex.FECHA_CHAVE) || a.testaTokenAtual(lex.FIM)) {
		instru := a.analisaInstrucao()
		if instru != nil {
			bloco.Instrucoes = append(bloco.Instrucoes, instru)
		}
		a.avancaToken()
	}
	a.profu--

	return bloco
}

func (a *Analisador) parseExpressaoPrefixo() (arv.Expressao, bool) {
	expr := &arv.ExpressaodePrefixo{Token: a.Atual, Operador: a.Atual.Valor, Profu: a.profu}
	var err bool
	a.avancaToken()
	expr.ExpDireita, err = a.parseExpressao(PREFIX)

	if !err {
		return nil, false
	}

	return expr, true
}

func (a *Analisador) parseExpressaoInfixo(expEsque arv.Expressao) (arv.Expressao, bool) {
	var err bool
	expre := &arv.ExpressaoInfixo{Token: a.Atual, Operador: a.Atual.Valor, ExpEsquerda: expEsque, Profu: a.profu}
	prioridade := a.prioridadeAtual()
	a.avancaToken()
	expre.ExpDireita, err = a.parseExpressao(prioridade)

	if !err {
		return nil, false
	}
	return expre, true
}

func (a *Analisador) parseExpressaoAtributo(exprEsq arv.Expressao) (arv.Expressao, bool) {
	expr := &arv.ExpressaoAtributo{Token: a.Atual, Profu: a.profu, Expres: exprEsq}

	if a.testaProximoToken(lex.IDENTIFICADOR) {
		a.avancaToken()
		expr.Atributo = a.Atual.Valor
		return expr, true
	} else {
		return a.parseExpressaoInfixo(exprEsq)
	}
}

func (a *Analisador) parseExpressaoLista() (arv.Expressao, bool) {
	expre := &arv.ExpressaoLista{Token: a.Atual, Profu: a.profu}
	filaExpre := aux.NewFilaExpressoes()

	a.avancaToken()
	if a.Atual.Tipo == lex.FECHA_COLCHETE {
		return expre, true
	}

	for {
		expreItem, ok := a.parseExpressao(MENOR)

		if !ok {
			return nil, false
		}

		filaExpre.AddItem(expreItem)

		if a.Proximo.Tipo == lex.VIRGULA {
			a.avancaToken()
			a.avancaToken()
			continue
		} else if a.Proximo.Tipo == lex.FECHA_COLCHETE {
			a.avancaToken()
			break
		} else {
			return nil, false
		}
	}

	expre.Expressoes = aux.GetLista(filaExpre)

	return expre, true
}

func (a *Analisador) parseExpressaoDict() (arv.Expressao, bool) {
	exprDict := &arv.ExpressaoDict{Token: a.Atual, Profu: a.profu}

	if a.testaProximoToken(lex.FECHA_CHAVE) {
		exprDict.Chaves = []arv.Expressao{}
		exprDict.Valores = []arv.Expressao{}
		a.avancaToken()
		return exprDict, true
	}

	var ok bool
	var chave arv.Expressao
	var value arv.Expressao

	a.avancaToken()
	filaChaves := aux.NewFilaExpressoes()
	filaValues := aux.NewFilaExpressoes()

	for {
		chave, ok = a.parseExpressao(MENOR)

		if !(ok && a.esperaToken(lex.SETA)) {
			return nil, false
		}

		a.avancaToken()
		value, ok = a.parseExpressao(MENOR)
		if !ok {
			return nil, false
		}

		filaChaves.AddItem(chave)
		filaValues.AddItem(value)

		if !a.testaProximoToken(lex.VIRGULA) {
			break
		}

		a.avancaToken()
		a.avancaToken()
	}

	if !a.esperaToken(lex.FECHA_CHAVE) {
		return nil, false
	}

	exprDict.Chaves = aux.GetLista(filaChaves)
	exprDict.Valores = aux.GetLista(filaValues)

	return exprDict, true
}

func (a *Analisador) parseExpressaoObjeto() (arv.Expressao, bool) {
	if !a.testaProximoToken(lex.ABRE_CHAVE) {
		return &arv.ChamadaObjeto{Token: a.Atual, Profu: a.profu}, true
	} else {
		exprObjeto := &arv.ExpressaoObjeto{Token: a.Atual, Profu: a.profu}
		exprObjeto.Atributos = make(map[string]arv.Expressao)

		a.avancaToken()

		for !a.testaProximoToken(lex.FECHA_CHAVE) {
			if a.esperaToken(lex.IDENTIFICADOR) {
				nomeAtrib := a.Atual.Valor

				if !a.esperaToken(lex.DOIS_PONTO) {
					return nil, false
				}

				a.avancaToken()
				expressao, ok := a.parseExpressao(MENOR)

				if !ok {
					return nil, false
				}

				exprObjeto.Atributos[nomeAtrib] = expressao

				if a.testaProximoToken(lex.VIRGULA) {
					a.avancaToken()
				}
			} else {
				return nil, false
			}
		}

		a.avancaToken()

		return exprObjeto, true
	}

}

func (a *Analisador) parseExpressaoClass() (arv.Expressao, bool) {
	//definindo os elementos do objeto de expressao
	expreClass := &arv.ExpressaoClass{
		Token:      a.Atual,
		AtribPub:   make(map[string]arv.Expressao),
		AtribPriv:  make(map[string]arv.Expressao),
		AtribPro:   make(map[string]arv.Expressao),
		AtribClass: make(map[string]arv.Expressao),
	}

	//verificamos se ha parenteses
	//se houver, significa que devemos declarar as super classes
	if a.testaProximoToken(lex.ABRE_PARENTESES) {
		a.avancaToken()

		//se for apenas um abre fecha chave, somente continua de forma normal
		if !a.testaProximoToken(lex.FECHA_PARENTESES) {
			filaSupers := a.getSupersClasse()

			if filaSupers == nil {
				return nil, false
			}

			expreClass.SuperClasses = aux.GetLista(filaSupers)
		}

		a.avancaToken()
	}

	if !a.esperaToken(lex.ABRE_CHAVE) {
		return nil, false
	}

	if a.testaProximoToken(lex.FECHA_CHAVE) {
		return expreClass, true
	}

	//e aqui que comeca a putaria!
	leveis := map[lex.TipoToken]int{
		lex.PUBLICO:   1,
		lex.PROTEGIDO: 2,
		lex.PRIVADO:   3,
		lex.CLASS:     4,
	}
	var nivelProtecao int = 1
	var nomeAtrib string
	var exprAtrib arv.Expressao
	var ok bool

	for {
		if a.testaProximoToken(lex.IDENTIFICADOR) {
			a.avancaToken()
		} else if a.testaProximoToken(lex.VIRGULA) {
			a.avancaToken()
			continue
		} else {
			nivelProtecao = leveis[a.Proximo.Tipo]
			if nivelProtecao < 1 || nivelProtecao > 4 {
				break
			}

			a.avancaToken()
			continue
		}

		nomeAtrib = a.Atual.Valor
		exprAtrib, ok = a.parseExpressao(MENOR)

		if !ok {
			return nil, false
		}

		switch nivelProtecao {
		case 1:
			expreClass.AtribPub[nomeAtrib] = exprAtrib
		case 2:
			expreClass.AtribPro[nomeAtrib] = exprAtrib
		case 3:
			expreClass.AtribPriv[nomeAtrib] = exprAtrib
		case 4:
			expreClass.AtribClass[nomeAtrib] = exprAtrib
		}
	}

	if a.esperaToken(lex.FECHA_CHAVE) {
		return expreClass, true
	}

	return nil, false
}

func (a *Analisador) getSupersClasse() *aux.FilaExpressoes {
	fila := aux.NewFilaExpressoes()

	for {
		a.avancaToken()
		expressao, ok := a.parseExpressao(MENOR)

		if !ok {
			return nil
		}

		fila.AddItem(expressao)

		if a.testaProximoToken(lex.VIRGULA) {
			a.avancaToken()
		} else if a.testaProximoToken(lex.FECHA_PARENTESES) {
			break
		} else {
			a.addErro(fmt.Sprintf("Token %s inesperado ou desconhecido", a.Proximo.Tipo))
			return nil
		}
	}

	return fila
}

func (a *Analisador) parseExpressaoAgrupada() (arv.Expressao, bool) {
	a.avancaToken()
	expr, err := a.parseExpressao(MENOR)

	if !err {
		return nil, false
	}

	if !a.esperaToken(lex.FECHA_PARENTESES) {
		return nil, false
	}

	return expr, true
}

func (a *Analisador) parseIdentificador() (arv.Expressao, bool) {
	return &arv.Identificador{Token: a.Atual, Nome: a.Atual.Valor, Profu: a.profu}, true
}

func (a *Analisador) parseLiteralInteiro() (arv.Expressao, bool) {
	expre := &arv.LiteralInt{Token: a.Atual, Profu: a.profu}
	expre.Valor, _ = strconv.ParseInt(a.Atual.Valor, 10, 64)

	return expre, true
}

func (a *Analisador) parseLiteralReal() (arv.Expressao, bool) {
	expre := &arv.LiteralReal{Token: a.Atual, Profu: a.profu}
	expre.Valor, _ = strconv.ParseFloat(a.Atual.Valor, 64)

	return expre, true
}

func (a *Analisador) parseBooleano() (arv.Expressao, bool) {
	return &arv.Booleano{Token: a.Atual, Valor: a.testaTokenAtual(lex.TRUE), Profu: a.profu}, true
}

func (a *Analisador) parseString() (arv.Expressao, bool) {
	return &arv.LiteralString{Token: a.Atual, Valor: a.Atual.Valor, Profu: a.profu}, true
}

func (a *Analisador) parseTipoNone() (arv.Expressao, bool) {
	return &arv.TipoNone{Token: a.Atual, Profu: a.profu}, true
}
