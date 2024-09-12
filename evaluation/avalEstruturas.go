package evaluation

import (
	obj "ChikenInterpreter/objetos"
	arv "ChikenInterpreter/parsing/arvore"
)

func avaliaIfElse(noIf *arv.ExpressaoIf, ambiente *obj.Ambiente) obj.ObjetoBase {
	var resultado obj.ObjetoBase

	condicao := Avaliar(noIf.Condicao, ambiente)
	if condicao.Tipo() == obj.ERRO {
		return condicao
	}

	entao := noIf.BlocoEntao
	senao := noIf.BlocoSenao

	if eVerdadeiro(condicao) {
		resultado = Avaliar(entao, ambiente)
	} else if senao != nil {
		resultado = Avaliar(senao, ambiente)
	} else {
		resultado = obj.OBJ_NONE
	}

	return resultado
}

func avaliaRepeat(noRepeat *arv.ExpressaoRepeat, ambiente *obj.Ambiente) obj.ObjetoBase {
	var item obj.ObjetoBase
	condicao1 := noRepeat.Condicao1
	condicao2 := noRepeat.Condicao2
	codigo := noRepeat.BlocoRepetir

	lista := make([]obj.ObjetoBase, 0)

	for {
		if !eVerdadeiro(Avaliar(condicao1, ambiente)) {
			break
		}

		item = avaliaInstrucoes(codigo.Instrucoes, ambiente)

		if item.Tipo() == obj.VALOR_RETORNO || item.Tipo() == obj.ERRO {
			return item
		} else if item == obj.OBJ_BREAK {
			break
		}

		if item != obj.OBJ_NONE && item != obj.OBJ_CONTINUE {
			lista = append(lista, item)
		}

		if !eVerdadeiro(Avaliar(condicao2, ambiente)) {
			break
		}
	}

	return &obj.ObjArray{ArrayList: lista, Capacidade: len(lista), Tamanho: len(lista)}
}
