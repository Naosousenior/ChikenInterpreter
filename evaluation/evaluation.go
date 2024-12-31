package evaluation

import (
	obj "ChikenInterpreter/objetos"
	arv "ChikenInterpreter/parsing/arvore"

	"fmt"
)

func FalsoHash(objeto obj.ObjetoBase) string {
	return string(objeto.Tipo()) + ": " + objeto.Inspecionar()
}

func status(tipo obj.TipoStatus,objeto obj.ObjetoBase) *obj.Status {
	return &obj.Status{Tipo: tipo,Resultado: objeto}
}

func AvaliaInstrucao(instrucao arv.Instrucao, ambiente *obj.Ambiente) *obj.Status {

	switch no := instrucao.(type) {
	case *arv.Programa:
		return avaliaPrograma(no.Instrucoes, ambiente)

	case *arv.InstrucaodeExpressao:
		switch expressao := no.Expressao.(type) {
		case *arv.ExpressaoIf:
			return avaliaIfElse(expressao,ambiente)
		case *arv.ExpressaoRepeat:
			return avaliaRepeat(expressao,ambiente)
		default:
			return status(obj.EXPRESSAO,AvaliaExpressao(expressao,ambiente))
		}

	case *arv.InstrucaoAtribuicao:
		novoValor := AvaliaExpressao(no.ExprValue, ambiente)

		if novoValor.Tipo() == obj.EXCECAO {
			return status(obj.ERROR,novoValor)
		}

		resultado := avaliaAtribuicao(no.Operador, no.ExprRecebe, novoValor, ambiente)

		if resultado.Tipo() == obj.EXCECAO {
			return status(obj.ERROR,resultado)
		}

		return status(obj.ATRIBUICAO,resultado)

	case *arv.ReturnInstrucao:
		res := AvaliaExpressao(no.Expre, ambiente)
		if res.Tipo() == obj.EXCECAO {
			return status(obj.ERROR,res)
		}
		return status(obj.RETURN,res)

	case *arv.InstrucaoBreak:
		return obj.BREAK_ST
	case *arv.InstrucaoContinue:
		return obj.CONTINUE_ST

	case *arv.ErrInstrucao:
		res := AvaliaExpressao(no.Expre, ambiente)

		return status(obj.ERROR,&obj.ObjExcessao{Mensagem: "Exeção: " + res.Inspecionar(), Objeto: res})
	
	case *arv.InstrucaoTryExcept:
		return avaliaTryExcept(no,ambiente);

	case *arv.VarInstrucao:
		for _, vardec := range no.Vars {
			valor := AvaliaExpressao(vardec.Expres, ambiente)
			if valor.Tipo() == obj.EXCECAO {
				return status(obj.ERROR,valor)
			}
			

			res := ambiente.CriaVar(vardec.Ident.Nome, valor)

			if res.Tipo() == obj.EXCECAO {
				return status(obj.ERROR,res)
			}
		}

		return obj.DECLARACAO_ST
	
	case *arv.DefInstrucao:
		nome := no.Ident.Nome
		valor := AvaliaExpressao(no.Expres,ambiente)

		if valor.Tipo() == obj.EXCECAO {
			return status(obj.ERROR,valor)
		}

		if ambiente.DefVar(nome,valor) {
			return obj.DEFINICAO_ST
		}

		return status(obj.ERROR,geraErro(fmt.Sprintf("Objeto %s já existe no escopo.",nome)))
	case *arv.InstrucaoIter:
		return avaliaIter(no, ambiente)

	case *arv.InstrucaoSwitch:
		return avaliaSwitch(no, ambiente)

	case *arv.BlocoInstrucao:
		return avaliaInstrucoes(no.Instrucoes, ambiente)
	}

	return status(obj.ERROR,geraErro("Instrucao desconhecida"))
}

func AvaliaExpressao(expressao arv.Expressao, ambiente *obj.Ambiente) obj.ObjetoBase {
	switch no := expressao.(type) {

	case *arv.ExpressaodePrefixo:
		exprDirei := AvaliaExpressao(no.ExpDireita, ambiente)
		if exprDirei.Tipo() == obj.EXCECAO {
			return exprDirei
		}
		return exprDirei.OpPrefixo(no.Operador)

	case *arv.ExpressaoAtributo:
		if _, ok := no.Expres.(*arv.ChamadaObjeto); ok {
			return ambiente.Objeto.Get(no.Atributo, ambiente)
		}

		objeto := AvaliaExpressao(no.Expres, ambiente)
		if objeto.Tipo() == obj.EXCECAO {
			return objeto
		}

		return objeto.GetPropriedade(no.Atributo)

	case *arv.ExpressaoIf:
		return avaliaIfElse(no, ambiente).Resultado

	case *arv.ExpressaoRepeat:
		return avaliaRepeat(no, ambiente).Resultado

	case *arv.ExpressaoFun:
		return &obj.ObjFuncao{Parametros: no.Parametros, BlocoInstrucoes: no.Bloco, Amb: ambiente}

	case *arv.ExpressaoClass:
		supers := make([]*obj.Classe, len(no.SuperClasses)+1)
		supers[len(supers)-1] = CLASSMAE

		for i, expr := range no.SuperClasses {
			resultado := AvaliaExpressao(expr, ambiente)
			if classe, ok := resultado.(*obj.Classe); ok {
				supers[i] = classe
			} else if resultado.Tipo() == obj.EXCECAO {
				return resultado
			} else {
				return geraErro(fmt.Sprintf("O objeto %s não é um objeto do tipo CLASS, e portanto não pode ser herdado", resultado.Inspecionar()))
			}
		}

		fmt.Println(ambiente)

		return avaliaClasse(no, supers, ambiente)
	case *arv.ExpressaoObjeto:

		retorno, erro := avaliaObject(no, CLASSMAE)

		if erro == nil {
			return retorno
		}

		return erro

	case *arv.CallFun:
		obj_fun := AvaliaExpressao(no.Funcao, ambiente)

		return avaliaChamada(no, obj_fun, ambiente)

	case *arv.ChamadaObjeto:
		if ambiente.Objeto == nil {
			return geraErro("Expressao 'object' fora de contexto")
		} else {
			return ambiente.Objeto
		}

	case *arv.ExpressaoInfixo:
		op := no.Operador
		esq := AvaliaExpressao(no.ExpEsquerda, ambiente)
		dir := AvaliaExpressao(no.ExpDireita, ambiente)

		if esq.Tipo() == obj.EXCECAO || dir.Tipo() == obj.EXCECAO {
			if esq.Tipo() == obj.EXCECAO {
				return esq
			} else {
				return dir
			}
		}

		return avaliaInfixo(op, esq, dir)

	case *arv.ExpressaoLista:
		valores := avaliaExpressoes(no.Expressoes, ambiente)
		if len(valores) > 0 && valores[0].Tipo() == obj.EXCECAO {
			return valores[0]
		}
		return &obj.ObjArray{ArrayList: valores, Capacidade: len(valores), Tamanho: len(valores)}

	case *arv.ExpressaoDict:
		return avaliaDict(no, ambiente)

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

	return geraErro(fmt.Sprintf("Tomar no boga meu irmao: %s", expressao))
}

func geraErro(msg string) *obj.ObjExcessao {
	return &obj.ObjExcessao{Mensagem: msg}
}

func avaliaPrograma(instrucoes []arv.Instrucao, ambiente *obj.Ambiente) *obj.Status {
	resultado := avaliaInstrucoes(instrucoes, ambiente)

	if resultado == nil {
		return &obj.Status{Tipo: obj.ERROR,Resultado: &obj.ObjExcessao{Mensagem: ""}}
	}

	return resultado
}

func avaliaInstrucoes(instrucoes []arv.Instrucao, ambiente *obj.Ambiente) *obj.Status {
	var resultado *obj.Status

	for _, instrucao := range instrucoes {
		resultado = AvaliaInstrucao(instrucao, ambiente)

		if resultado == nil {
			fmt.Println("Cara, não é pra retornar nil não")
			continue
		} else if resultado.Tipo == obj.RETURN {
			return resultado
		} else if resultado.Tipo == obj.ERROR {
			fmt.Printf("Linha: %d\n", instrucao.GetTokenNo().Pos.Linha)
			break
		} else if resultado == obj.BREAK_ST || resultado == obj.CONTINUE_ST {
			break
		}
	}

	return resultado
}

func avaliaExpressoes(expressoes []arv.Expressao, ambiente *obj.Ambiente) []obj.ObjetoBase {
	resultado := make([]obj.ObjetoBase, len(expressoes))

	for i, exp := range expressoes {
		av := AvaliaExpressao(exp, ambiente)

		if av.Tipo() == obj.EXCECAO {
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
