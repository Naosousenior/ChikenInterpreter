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
			if novoValor.Tipo() == obj.ERRO {
				return novoValor
			}
			ambiente.SetVar(expr.Nome, novoValor)

			return obj.OBJ_NONE
		} else {
			return geraErro(fmt.Sprintf("Variavel %s nao declarada", expr.Nome))
		}
	case *arv.ExpressaoAtributo:
		var ok bool

		objeto := Avaliar(expr.Expres, ambiente)
		if objeto.Tipo() == obj.ERRO {
			return objeto
		}

		var velhoValor obj.ObjetoBase

		if _,ok = expr.Expres.(*arv.ChamadaObjeto); ok {
			velhoValor = ambiente.Objeto.Get(expr.Atributo,ambiente)
		} else {
			velhoValor = objeto.GetPropriedade(expr.Atributo)
		}

		novoValor = avaliaOperadorAtribuicao(operador, novoValor, velhoValor)

		if novoValor.Tipo() == obj.ERRO {
			return novoValor
		}

		if ok {
			return ambiente.Objeto.Set(expr.Atributo,novoValor,ambiente)
		}

		return objeto.SetPropriedade(expr.Atributo, novoValor)

	case *arv.ExpressaoInfixo:
		if expr.Operador == "[" {
			obj_aux := Avaliar(expr.ExpEsquerda, ambiente)
			objeto, ok := obj_aux.(obj.ObjetoIndexavel)
			if !ok {
				return geraErro(fmt.Sprintf("Objeto %s não indexável", obj_aux.Inspecionar()))
			}

			index := Avaliar(expr.ExpDireita, ambiente)

			if objeto.Tipo() == obj.ERRO {
				return objeto
			}

			if index.Tipo() == obj.ERRO {
				return objeto
			}

			velhoValor := objeto.OpInfixo("[", index)
			novoValor = avaliaOperadorAtribuicao(operador, novoValor, velhoValor)

			if novoValor.Tipo() == obj.ERRO {
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
		if velhoValor.Tipo() == obj.ERRO {
			return velhoValor
		}

		if novoValor.Tipo() == velhoValor.Tipo() {
			return novoValor
		} else {
			return geraErro("O tipo recebido e incompatível com o tipo atribuido")
		}
	} else {

		if velhoValor.Tipo() == obj.ERRO {
			return velhoValor
		}
		novoValor = velhoValor.OpInfixo(op[0:1], novoValor)
		return novoValor
	}
}

func avaliaIter(noIter *arv.InstrucaoIter, ambiente *obj.Ambiente) obj.ObjetoBase {
	var ok bool
	var lista obj.ObjetoIndexavel
	iterador := noIter.Iterador.Nome
	blocoCodigo := noIter.BlocoCodigo
	objeto := Avaliar(noIter.ExpressaoLista, ambiente)

	if lista, ok = objeto.(obj.ObjetoIndexavel); !ok {
		return geraErro("Instrução de iteração precisa de um objeto iteravel")
	}

	count := 0
	novoAmbiente := obj.NewAmbienteInterno(ambiente)
	if objeto = novoAmbiente.CriaVar(iterador, obj.OBJ_NONE); objeto.Tipo() == obj.ERRO {
		return objeto
	}

	var valorAtual obj.ObjetoBase = obj.OBJ_NONE

	for {
		valorAtual = lista.Iterar(count)
		if valorAtual == obj.OBJ_BREAK {
			break
		}

		novoAmbiente.SetVar(iterador, valorAtual)
		count++

		objeto = avaliaInstrucoes(blocoCodigo.Instrucoes, novoAmbiente)

		if objeto.Tipo() == obj.ERRO || objeto.Tipo() == obj.VALOR_RETORNO {
			return objeto
		}

		if objeto == obj.OBJ_BREAK {
			break
		}
	}

	return obj.OBJ_NONE
}

func avaliaSwitch(noSwitch *arv.InstrucaoSwitch, ambiente *obj.Ambiente) obj.ObjetoBase {
	casos := noSwitch.Cases
	var valorCaso obj.ObjetoBase
	var hash string

	valorCambio := Avaliar(noSwitch.ExpreTeste,ambiente)

	if valorCambio.Tipo() == obj.ERRO {
		return valorCambio
	}

	mapaDeCasos := make(map[string]*arv.BlocoInstrucao)

	for i, caso := range casos {
		valorCaso = Avaliar(caso.ExpreCase,ambiente)

		if valorCaso.Tipo() == obj.ERRO {
			return valorCaso
		}

		hash = FalsoHash(valorCaso)
		mapaDeCasos[hash] = casos[i].Codigo
	}

	hash = FalsoHash(valorCambio)
	blocoCodigo,ok := mapaDeCasos[hash]

	if !ok {
		if noSwitch.BlocoDefault == nil {
			return obj.OBJ_NONE
		}

		return avaliaInstrucoes(noSwitch.BlocoDefault.Instrucoes,ambiente)
	}

	return avaliaInstrucoes(blocoCodigo.Instrucoes,ambiente)
}