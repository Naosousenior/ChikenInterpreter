package objetos

import (
	"ChikenInterpreter/lexing"
	arv "ChikenInterpreter/parsing/arvore"
	"fmt"
	"math"
	"strconv"
	"strings"
)

func geraErro(erro string) ObjetoBase {
	return &ObjExcessao{Mensagem: erro}
}

//comecamos implementando os numeros

func calculaInfixNum(operador string, vEsq, vDir ObjetoBase) ObjetoBase {
	var v1, v2, res float64

	v1, _ = strconv.ParseFloat(vEsq.Inspecionar(), 64)
	v2, _ = strconv.ParseFloat(vDir.Inspecionar(), 64)

	switch operador {
	case "+":
		res = v1 + v2
	case "-":
		res = v1 - v2
	case "*":
		res = v1 * v2
	case "/":

		if v2 == 0 {
			return geraErro("Nenhum valor pode ser dividido por 0. ")
		}
		res = v1 / v2
	case "**":
		res = math.Pow(v1, v2)
	case "%":
		res = math.Mod(v1, v2)
	case "<":
		if v1 < v2 {
			return OBJ_TRUE
		}
		return OBJ_FALSE
	case "<=":
		if v1 <= v2 {
			return OBJ_TRUE
		}
		return OBJ_FALSE
	case ">":
		if v1 > v2 {
			return OBJ_TRUE
		}

		return OBJ_FALSE
	case ">=":
		if v1 >= v2 {
			return OBJ_TRUE
		}

		return OBJ_FALSE
	case "==":
		if v1 == v2 {
			return OBJ_TRUE
		}

		return OBJ_FALSE
	case "!=":
		if v1 != v2 {
			return OBJ_TRUE
		}

		return OBJ_FALSE
	default:
		return geraErro(fmt.Sprintf("O operador %s nao pode ser usado entre dois numeros. ", operador))
	}

	inteiro, fracao := math.Modf(res)

	if fracao == 0 {
		return &ObjInteiro{Valor: int(inteiro)}
	} else {
		return &ObjReal{Valor: res}
	}
}

type ObjInteiro struct {
	Valor int
}

func (oi *ObjInteiro) Tipo() TipoObjeto    { return INTEIRO }
func (oi *ObjInteiro) Inspecionar() string { return fmt.Sprintf("%d", oi.Valor) }
func (oi *ObjInteiro) OpInfixo(op string, dir ObjetoBase) ObjetoBase {
	if dir.Tipo() == REAL || dir.Tipo() == INTEIRO {
		return calculaInfixNum(op, oi, dir)
	}

	return geraErro(fmt.Sprintf("Objeto %s incompativel com %s na operacao %s", oi.Tipo(), dir.Tipo(), op))
}
func (oi *ObjInteiro) OpPrefixo(op string) ObjetoBase {
	if op == lexing.SUB {
		return &ObjInteiro{Valor: -oi.Valor}
	} else {
		return geraErro(fmt.Sprintf("Tipo INTEIRO incompativel com operacao %s", op))
	}
}
func (oe *ObjInteiro) GetPropriedade(propri string) ObjetoBase {
	return geraErro("Objeto INTEIRO nao possui propriedades")
}
func (oe *ObjInteiro) SetPropriedade(propri string, obj ObjetoBase) ObjetoBase {
	return geraErro("Objeto INTEIRO nao possui propriedades")
}

type ObjReal struct {
	Valor float64
}

func (or *ObjReal) Tipo() TipoObjeto { return REAL }
func (or *ObjReal) Inspecionar() string {
	return fmt.Sprintf("%f", or.Valor)
}

func (or *ObjReal) OpInfixo(op string, dir ObjetoBase) ObjetoBase {
	if dir.Tipo() == REAL || dir.Tipo() == INTEIRO {
		return calculaInfixNum(op, or, dir)
	}

	return geraErro(fmt.Sprintf("Objeto %s incompativel com %s na operacao %s", or.Tipo(), dir.Tipo(), op))
}

func (or *ObjReal) OpPrefixo(op string) ObjetoBase {
	if op == lexing.SUB {
		return &ObjReal{Valor: -or.Valor}
	} else {
		return geraErro(fmt.Sprintf("Operacao %s incompativel com tipo REAL", op))
	}
}
func (oe *ObjReal) GetPropriedade(propri string) ObjetoBase {
	return geraErro("Objeto REAL nao possui propriedades")
}
func (oe *ObjReal) SetPropriedade(propri string, obj ObjetoBase) ObjetoBase {
	return geraErro("Objeto REAL nao possui propriedades")
}

// implementacao de strings

type ObjTexto struct {
	Valor string
}

func (ot *ObjTexto) Tipo() TipoObjeto    { return TEXTO }
func (ot *ObjTexto) Inspecionar() string { return ot.Valor }
func (ot *ObjTexto) OpPrefixo(op string) ObjetoBase {
	return geraErro("Strings não suportam prefixos")
}
func (ot *ObjTexto) OpInfixo(op string, direita ObjetoBase) ObjetoBase {
	switch op {
	case "+":
		return &ObjTexto{Valor: ot.Valor + direita.Inspecionar()}

	case "[":
		valor, ok := direita.(*ObjInteiro)

		if !ok {
			return geraErro("O indexador de strings deve ser um inteiro")
		}

		pos := valor.Valor

		return &ObjTexto{Valor: fmt.Sprintf("%c", ot.Valor[pos])}
	}

	return geraErro(fmt.Sprintf("Opercao %s nao suportada por strings", op))
}
func (ot *ObjTexto) GetPropriedade(propri string) ObjetoBase {
	return geraErro(fmt.Sprintf("Propriedade %s não encontrada", propri))
}
func (ot *ObjTexto) SetPropriedade(propri string, obj ObjetoBase) ObjetoBase {
	return geraErro(fmt.Sprintf("Propriedade %s não encontrada", propri))
}
func (ot *ObjTexto) SetIndex(index ObjetoBase, valor ObjetoBase) ObjetoBase {
	return geraErro("Nao e possivel alterar uma posicao de uma string")
}

//agora implementacao dos booleanos

func calculaInfixBool(operador string, vEsq, vDir ObjetoBase) ObjetoBase {
	v1 := OBJ_FALSE
	if vEsq.Inspecionar() == "TRUE" {
		v1 = OBJ_TRUE
	}

	v2 := OBJ_TRUE
	if vDir.Inspecionar() == "FALSE" {
		v2 = OBJ_FALSE
	}

	switch operador {
	case "&":
		if v1.Valor && v2.Valor {
			return OBJ_TRUE
		}

		return OBJ_FALSE
	case "|":
		if v1.Valor || v2.Valor {
			return OBJ_TRUE
		}

		return OBJ_FALSE
	case "||":
		if v1.Valor == v2.Valor {
			return OBJ_FALSE
		}

		return OBJ_TRUE

	case "==":
		if v1.Valor == v2.Valor {
			return OBJ_TRUE
		}

		return OBJ_FALSE
	case "!=":
		if v1.Valor != v2.Valor {
			return OBJ_TRUE
		}

		return OBJ_FALSE
	default:
		return geraErro(fmt.Sprintf("O operador %s nao pode ser usado entre dois booleanos. ", operador))
	}
}

type ObjBool struct {
	Valor bool
}

func (ob *ObjBool) Tipo() TipoObjeto { return BOOLEANO }
func (ob *ObjBool) Inspecionar() string {
	if ob.Valor {
		return "TRUE"
	}

	return "FALSE"
}
func (ob *ObjBool) OpInfixo(op string, dir ObjetoBase) ObjetoBase {
	if ob.Tipo() == BOOLEANO {
		return calculaInfixBool(op, ob, dir)
	}

	return geraErro(fmt.Sprintf("Objeto %s incompativel com %s na operacao %s", ob.Tipo(), dir.Tipo(), op))
}
func (ob *ObjBool) OpPrefixo(op string) ObjetoBase {
	if op == lexing.NEGACAO {
		if ob.Valor {
			return OBJ_FALSE
		} else {
			return OBJ_TRUE
		}
	} else {
		return geraErro(fmt.Sprintf("Objeto BOOL incompativel com operacao %s", op))
	}
}
func (oe *ObjBool) GetPropriedade(propri string) ObjetoBase {
	return geraErro("Objeto BOOL nao possui propriedades")
}
func (oe *ObjBool) SetPropriedade(propri string, obj ObjetoBase) ObjetoBase {
	return geraErro("Objeto BOOL nao possui propriedades")
}

type ObjNone struct {
}

func (on *ObjNone) Tipo() TipoObjeto    { return NONE }
func (on *ObjNone) Inspecionar() string { return "NONE" }
func (on *ObjNone) OpInfixo(op string, dir ObjetoBase) ObjetoBase {
	if op == "==" || op == "!=" {
		if op == "==" && dir == OBJ_NONE {
			return OBJ_TRUE
		}

		if op == "!=" && dir != OBJ_NONE {
			return OBJ_TRUE
		}

		return OBJ_FALSE
	}
	return geraErro("Tipo NONE incompativel com qualquer operacao")
}
func (on *ObjNone) OpPrefixo(op string) ObjetoBase {
	return geraErro("Tipo NONE incompativel com qualquer operacao")
}
func (oe *ObjNone) GetPropriedade(propri string) ObjetoBase {
	return geraErro("Tipo NONE nao possui propriedades")
}
func (oe *ObjNone) SetPropriedade(propri string, obj ObjetoBase) ObjetoBase {
	return geraErro("Tipo NONE nao possui propriedades")
}

//implementando dois objetos usados internamente

type ObjExcessao struct {
	Mensagem string
	Objeto   ObjetoBase
}

func (oe *ObjExcessao) Tipo() TipoObjeto    { return EXCECAO }
func (oe *ObjExcessao) Inspecionar() string {
	if oe.Objeto == nil {
		return oe.Mensagem
	}

	return fmt.Sprintf("%s: %s",oe.Mensagem,oe.Objeto)
}
func (oe *ObjExcessao) OpInfixo(operador string, direita ObjetoBase) ObjetoBase {
	return geraErro("Objetos do tipo ERRO nao realizam operacoes")
}
func (oe *ObjExcessao) OpPrefixo(op string) ObjetoBase {
	return geraErro("Objetos do tipo ERRO nao realizam operacoes")
}
func (oe *ObjExcessao) GetPropriedade(propri string) ObjetoBase {
	return geraErro("")
}
func (oe *ObjExcessao) SetPropriedade(propri string, obj ObjetoBase) ObjetoBase {
	return geraErro("")
}

type ObjFuncao struct {
	Amb             *Ambiente
	Parametros      []*arv.Identificador
	BlocoInstrucoes *arv.BlocoInstrucao
}

func (f *ObjFuncao) Tipo() TipoObjeto { return FUNCAO_OBJ }
func (f *ObjFuncao) Inspecionar() string {
	partes := make([]string, len(f.Parametros)+2)

	partes[0] = "function ("

	for i, v := range f.Parametros {
		partes[i+1] = v.Nome + ", "
	}

	partes[len(partes)-1] = ")"

	return strings.Join(partes, "")
}

func (f *ObjFuncao) OpInfixo(op string, dir ObjetoBase) ObjetoBase {
	return geraErro("Funcoes nao suportam operacoes")
}

func (f *ObjFuncao) OpPrefixo(op string) ObjetoBase {
	return geraErro("Funcoes nao suportam operacoes")
}
func (f *ObjFuncao) GetPropriedade(propri string) ObjetoBase {
	return geraErro("Funcoes nao possuem propriedades")
}
func (f *ObjFuncao) SetPropriedade(propri string, valor ObjetoBase) ObjetoBase {
	return geraErro("Funcoes nao possuem propriedades")
}
