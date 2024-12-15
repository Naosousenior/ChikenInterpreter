package evaluation

import (
	obj "ChikenInterpreter/objetos"
	arv "ChikenInterpreter/parsing/arvore"

	"fmt"
)

func FalsoHash(objeto obj.ObjetoBase) string {
	return string(objeto.Tipo()) +": "+objeto.Inspecionar()
}

func Avaliar(no arv.No, ambiente *obj.Ambiente) obj.ObjetoBase {
	switch no := no.(type) {
	case *arv.Programa:
		return avaliaPrograma(no.Instrucoes, ambiente)

	case *arv.InstrucaodeExpressao:
		return Avaliar(no.Expressao, ambiente)

	case *arv.InstrucaoAtribuicao:
		novoValor := Avaliar(no.ExprValue, ambiente)

		if novoValor.Tipo() == obj.ERRO {
			return novoValor
		}

		return avaliaAtribuicao(no.Operador, no.ExprRecebe, novoValor, ambiente)

	case *arv.ReturnInstrucao:
		res := Avaliar(no.Expre, ambiente)
		return &obj.ObjReturn{Valor: res}

	case *arv.InstrucaoBreak:
		return obj.OBJ_BREAK
	case *arv.InstrucaoContinue:
		return obj.OBJ_CONTINUE

	case *arv.ErrInstrucao:
		res := Avaliar(no.Expre, ambiente)

		return &obj.ObjErro{Mensagem: "Exeção: " + res.Inspecionar(), Objeto: res}

	case *arv.VarInstrucao:
		for _, vardec := range no.Vars {
			valor := Avaliar(vardec.Expres, ambiente)
			if valor.Tipo() == obj.ERRO {
				return valor
			}

			if aux, ok := valor.(*obj.ObjReturn); ok {
				valor = aux.Valor
			}

			res := ambiente.CriaVar(vardec.Ident.Nome, valor)

			if res.Tipo() == obj.ERRO {
				return res
			}
		}

		return obj.OBJ_NONE

	case *arv.InstrucaoIter:
		return avaliaIter(no, ambiente)

	case *arv.InstrucaoSwitch:
		return avaliaSwitch(no,ambiente)

	case *arv.ExpressaodePrefixo:
		exprDirei := Avaliar(no.ExpDireita, ambiente)
		if exprDirei.Tipo() == obj.ERRO {
			return exprDirei
		}
		return exprDirei.OpPrefixo(no.Operador)

	case *arv.ExpressaoAtributo:
		if _,ok := no.Expres.(*arv.ChamadaObjeto); ok{
			return ambiente.Objeto.Get(no.Atributo,ambiente)
		}

		objeto := Avaliar(no.Expres, ambiente)
		if objeto.Tipo() == obj.ERRO {
			return objeto
		}

		return objeto.GetPropriedade(no.Atributo)

	case *arv.BlocoInstrucao:
		return avaliaInstrucoes(no.Instrucoes, ambiente)

	case *arv.ExpressaoIf:
		return avaliaIfElse(no, ambiente)

	case *arv.ExpressaoRepeat:
		return avaliaRepeat(no, ambiente)

	case *arv.ExpressaoFun:
		return &obj.ObjFuncao{Parametros: no.Parametros, BlocoInstrucoes: no.Bloco, Amb: ambiente}

	case *arv.ExpressaoClass:
		supers := make([]*obj.Classe, len(no.SuperClasses)+1)
		supers[len(supers)-1] = CLASSMAE

		for i,expr := range no.SuperClasses {
			resultado := Avaliar(expr,ambiente)
			if classe,ok := resultado.(*obj.Classe);ok {
				supers[i] = classe
			} else if resultado.Tipo() == obj.ERRO {
				return resultado
			} else {
				return geraErro(fmt.Sprintf("O objeto %s não é um objeto do tipo CLASS, e portanto não pode ser herdado",resultado.Inspecionar()))
			}
		}

		return avaliaClasse(no,supers,ambiente)
	case *arv.ExpressaoObjeto:

		retorno,erro := avaliaObject(no,CLASSMAE)

		if erro == nil {
			return retorno
		}

		return erro

	case *arv.CallFun:
		obj_fun := Avaliar(no.Funcao, ambiente)

		return avaliaChamada(no,obj_fun,ambiente)

	case *arv.ChamadaObjeto:
		if ambiente.Objeto == nil {
			return geraErro("Expressao 'object' fora de contexto")
		} else {
			return ambiente.Objeto
		}

	case *arv.ExpressaoInfixo:
		op := no.Operador
		esq := Avaliar(no.ExpEsquerda, ambiente)
		dir := Avaliar(no.ExpDireita, ambiente)

		if esq.Tipo() == obj.ERRO || dir.Tipo() == obj.ERRO {
			if esq.Tipo() == obj.ERRO {
				return esq
			} else {
				return dir
			}
		}

		return avaliaInfixo(op, esq, dir)

	case *arv.ExpressaoLista:
		valores := avaliaExpressoes(no.Expressoes, ambiente)
		if len(valores) > 0 && valores[0].Tipo() == obj.ERRO {
			return valores[0]
		}
		return &obj.ObjArray{ArrayList: valores, Capacidade: len(valores), Tamanho: len(valores)}

	case *arv.ExpressaoDict:
		return avaliaDict(no,ambiente)

	case *arv.LiteralInt:
		return &obj.ObjInteiro{Valor: int(no.Valor)}
	case *arv.Booleano:
		if no.Valor {
			return obj.OBJ_TRUE
		} else {
			return obj.OBJ_FALSE
		}
	case *arv.LiteralReal:
		return &obj.ObjReal{Valor: no.Valor}

	case *arv.LiteralString:
		return &obj.ObjTexto{Valor: no.Valor}

	case *arv.Identificador:
		return avaliaIdentificador(no.Nome, ambiente)
	case *arv.TipoNone:
		return obj.OBJ_NONE
	}

	return geraErro("Tomar no boga meu irmao")
}

func geraErro(msg string) *obj.ObjErro {
	return &obj.ObjErro{Mensagem: msg}
}

func avaliaPrograma(instrucoes []arv.Instrucao, ambiente *obj.Ambiente) obj.ObjetoBase {
	resultado := avaliaInstrucoes(instrucoes, ambiente)

	if retorno, ok := resultado.(*obj.ObjReturn); ok {
		return retorno.Valor
	} else if retorno, ok := resultado.(*obj.ObjInstrucao); ok {
		return geraErro(fmt.Sprintf("Instrucao %s fora de contexto",retorno.Inspecionar()))
	}

	if resultado == nil {
		return &obj.ObjErro{Mensagem: ""}
	}

	return resultado
}

func avaliaInstrucoes(instrucoes []arv.Instrucao, ambiente *obj.Ambiente) obj.ObjetoBase {
	var resultado obj.ObjetoBase

	for _, instrucao := range instrucoes {
		resultado = Avaliar(instrucao, ambiente)

		if resultado == nil {
			continue
		} else if resultado.Tipo() == obj.VALOR_RETORNO {
			return resultado
		} else if resultado.Tipo() == obj.ERRO  {
			fmt.Printf("Linha: %d\n", instrucao.GetTokenNo().Pos.Linha)
			break
		} else if resultado == obj.OBJ_BREAK || resultado == obj.OBJ_CONTINUE {
			break
		}
	}

	return resultado
}

func avaliaExpressoes(expressoes []arv.Expressao, ambiente *obj.Ambiente) []obj.ObjetoBase {
	resultado := make([]obj.ObjetoBase, len(expressoes))

	for i, exp := range expressoes {
		av := Avaliar(exp, ambiente)

		if av.Tipo() == obj.ERRO {
			return []obj.ObjetoBase{av}
		}

		resultado[i] = av
	}

	return resultado
}

func eVerdadeiro(objeto obj.ObjetoBase) bool {
	switch obj := objeto.(type) {
	case *obj.ObjBool:
		if obj.Valor {
			return true
		}
		return false

	case *obj.ObjInteiro:
		if obj.Valor > 0 {
			return true
		}
		return false

	case *obj.ObjReal:
		if obj.Valor > 0 {
			return true
		}

		return false

	case *obj.ObjTexto:
		if obj.Valor != "" {
			return true
		}

		return false
	}

	return false
}
