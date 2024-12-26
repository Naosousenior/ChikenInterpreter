package evaluation

import (
	obj "ChikenInterpreter/objetos"
	arv "ChikenInterpreter/parsing/arvore"
)

func avaliaIfElse(noIf *arv.ExpressaoIf, ambiente *obj.Ambiente) *obj.Status {
	var resultado *obj.Status

	condicao := AvaliaExpressao(noIf.Condicao, ambiente)
	if condicao.Tipo() == obj.EXCECAO {
		return status(obj.ERROR,condicao)
	}

	entao := noIf.BlocoEntao
	senao := noIf.BlocoSenao

	if eVerdadeiro(condicao) {
		resultado = avaliaInstrucoes(entao.Instrucoes, ambiente)
	} else if senao != nil {
		resultado = avaliaInstrucoes(senao.Instrucoes, ambiente)
	} else {
		resultado = status(obj.EXPRESSAO,obj.OBJ_NONE)
	}

	return resultado
}

func avaliaRepeat(noRepeat *arv.ExpressaoRepeat, ambiente *obj.Ambiente) *obj.Status {
	var item *obj.Status
	condicao1 := noRepeat.Condicao1
	condicao2 := noRepeat.Condicao2
	codigo := noRepeat.BlocoRepetir

	lista := make([]obj.ObjetoBase, 0)

	for {
		if !eVerdadeiro(AvaliaExpressao(condicao1, ambiente)) {
			break
		}

		item = avaliaInstrucoes(codigo.Instrucoes, ambiente)

		if item.Tipo == obj.RETURN || item.Tipo == obj.ERROR {
			return item
		} else if item == obj.BREAK_ST {
			break
		}

		if item.Resultado != obj.OBJ_NONE && item != obj.CONTINUE_ST {
			lista = append(lista, item.Resultado)
		}

		if !eVerdadeiro(AvaliaExpressao(condicao2, ambiente)) {
			break
		}
	}

	return status(obj.EXPRESSAO,&obj.ObjArray{ArrayList: lista, Capacidade: len(lista), Tamanho: len(lista)})
}
