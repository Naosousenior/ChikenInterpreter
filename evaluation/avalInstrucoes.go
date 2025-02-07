package evaluation

import (
	"ChikenInterpreter/lexing"
	obj "ChikenInterpreter/objetos"
	arv "ChikenInterpreter/parsing/arvore"

	"fmt"
)

func avaliaAtribuicao(operador string, receptor arv.Expressao, novoValor obj.ObjetoBase, ambiente *obj.Ambiente) obj.ObjetoBase {

	switch expr := receptor.(type) {
	case *arv.Identificador:
		velhoValor, ok := ambiente.GetVar(expr.Nome)

		if ok {
			novoValor = avaliaOperadorAtribuicao(operador, novoValor, velhoValor)
			if novoValor.Tipo() == obj.EXCECAO {
				return novoValor
			}
			if !ambiente.SetVar(expr.Nome, novoValor) {
				return geraErro(fmt.Sprintf("Variável %s não pode ser reatribuida",expr.Nome))
			}

			return obj.OBJ_NONE
		} else {
			return geraErro(fmt.Sprintf("Variavel %s nao declarada.", expr.Nome))
		}
	case *arv.ExpressaoAtributo:
		var ok bool

		objeto := AvaliaExpressao(expr.Expres, ambiente)
		if objeto.Tipo() == obj.EXCECAO {
			return objeto
		}

		var velhoValor obj.ObjetoBase

		if _,ok = expr.Expres.(*arv.ChamadaObjeto); ok {
			velhoValor = ambiente.Objeto.Get(expr.Atributo,ambiente)
		} else {
			velhoValor = objeto.GetPropriedade(expr.Atributo)
		}

		novoValor = avaliaOperadorAtribuicao(operador, novoValor, velhoValor)

		if novoValor.Tipo() == obj.EXCECAO {
			return novoValor
		}

		if ok {
			return ambiente.Objeto.Set(expr.Atributo,novoValor,ambiente)
		}

		return objeto.SetPropriedade(expr.Atributo, novoValor)

	case *arv.ExpressaoInfixo:
		if expr.Operador == "[" {
			obj_aux := AvaliaExpressao(expr.ExpEsquerda, ambiente)
			objeto, ok := obj_aux.(obj.ObjetoIndexavel)
			if !ok {
				return geraErro(fmt.Sprintf("Objeto %s não indexável", obj_aux.Inspecionar()))
			}

			index := AvaliaExpressao(expr.ExpDireita, ambiente)

			if objeto.Tipo() == obj.EXCECAO {
				return objeto
			}

			if index.Tipo() == obj.EXCECAO {
				return objeto
			}

			velhoValor := objeto.OpInfixo("[", index)
			novoValor = avaliaOperadorAtribuicao(operador, novoValor, velhoValor)

			if novoValor.Tipo() == obj.EXCECAO {
				return novoValor
			}

			return objeto.SetIndex(index, novoValor)
		}

		return geraErro("Expressao " + expr.GetInformacao() + " nao atribuivel")

	default:
		return geraErro("Expressao " + expr.GetInformacao() + " nao atribuivel")
	}
}

func avaliaOperadorAtribuicao(op string, novoValor, velhoValor obj.ObjetoBase) obj.ObjetoBase {
	if op == lexing.RECEBE {
		return novoValor
	} else if op == lexing.TIPO_RECEBE {
		if velhoValor.Tipo() == obj.EXCECAO {
			return velhoValor
		}

		if novoValor.Tipo() == velhoValor.Tipo() {
			return novoValor
		} else {
			return geraErro("O tipo recebido e incompatível com o tipo atribuido")
		}
	} else {

		if velhoValor.Tipo() == obj.EXCECAO {
			return velhoValor
		}
		novoValor = velhoValor.OpInfixo(op[0:1], novoValor)
		return novoValor
	}
}

func avaliaIter(noIter *arv.InstrucaoIter, ambiente *obj.Ambiente) *obj.Status {
	objeto := AvaliaExpressao(noIter.ExpressaoLista,ambiente)

	if objeto.Tipo() == obj.EXCECAO {
		return status(obj.ERROR,objeto)
	}

	if iterado,ok := objeto.(obj.ObjetoIndexavel); ok {
		iterador := noIter.Iterador.Nome
		ambInterno := obj.NewAmbienteInterno(ambiente)
		var resultado *obj.Status
		for valor := range iterado.Iterar() {
			ambInterno.AddArgs(iterador,valor)

			resultado = avaliaInstrucoes(noIter.BlocoCodigo.Instrucoes,ambInterno)

			if resultado == obj.BREAK_ST {
				break
			} else if resultado.Tipo == obj.RETURN || resultado.Tipo == obj.ERROR {
				return resultado
			}
		}

		return obj.ITER_ST
	}

	return status(obj.ERROR,geraErro(fmt.Sprintf("Objeto %s não pode ser iterado",objeto.Inspecionar())))
}

func avaliaSwitch(noSwitch *arv.InstrucaoSwitch, ambiente *obj.Ambiente) *obj.Status {
	casos := noSwitch.Cases
	var valorCaso obj.ObjetoBase
	var hash string

	valorCambio := AvaliaExpressao(noSwitch.ExpreTeste,ambiente)

	if valorCambio.Tipo() == obj.EXCECAO {
		return status(obj.ERROR,valorCambio)
	}

	mapaDeCasos := make(map[string]*arv.BlocoInstrucao)

	for i, caso := range casos {
		valorCaso = AvaliaExpressao(caso.ExpreCase,ambiente)

		if valorCaso.Tipo() == obj.EXCECAO {
			return status(obj.ERROR,valorCaso)
		}

		hash = FalsoHash(valorCaso)
		mapaDeCasos[hash] = casos[i].Codigo
	}

	hash = FalsoHash(valorCambio)
	blocoCodigo,ok := mapaDeCasos[hash]

	if !ok {
		if noSwitch.BlocoDefault == nil {
			return obj.SWITCH_ST
		}

		return avaliaInstrucoes(noSwitch.BlocoDefault.Instrucoes,ambiente)
	}

	return avaliaInstrucoes(blocoCodigo.Instrucoes,ambiente)
}

func avaliaTryExcept(noTry *arv.InstrucaoTryExcept,ambiente *obj.Ambiente) *obj.Status {
	novoAmb := obj.NewAmbienteInterno(ambiente)

	stts_result := avaliaInstrucoes(noTry.BlocoTry.Instrucoes,novoAmb)

	if stts_result.Tipo == obj.ERROR {
		obj_exce := stts_result.Resultado
		ambExcept := obj.NewAmbienteInterno(ambiente)
		ambExcept.AddArgs(noTry.ExcessaoVar,obj_exce)

		return avaliaInstrucoes(noTry.BlocoExcept.Instrucoes,ambExcept)
	}

	return stts_result
}