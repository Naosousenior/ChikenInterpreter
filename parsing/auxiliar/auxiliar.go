package auxiliar

import "ChikenInterpreter/parsing/arvore"

type FilaExpressoes struct {
	Tamanho  int
	primeiro *ItemFila
	ultimo   *ItemFila
}

type ItemFila struct {
	item    arvore.Expressao
	proximo *ItemFila
}

func NewFilaExpressoes() *FilaExpressoes {
	fila := &FilaExpressoes{Tamanho: 0}
	return fila
}

func (fl *FilaExpressoes) AddItem(expr arvore.Expressao) {
	novo := &ItemFila{item: expr}
	if fl.Tamanho == 0 {
		fl.primeiro = novo
		fl.ultimo = novo
		fl.Tamanho++
		return
	}
	fl.ultimo.proximo = novo
	fl.ultimo = novo
	fl.Tamanho++
}
func (fl *FilaExpressoes) GetItem() arvore.Expressao {
	expressao := fl.primeiro.item
	fl.primeiro = fl.primeiro.proximo
	//fl.Tamanho--
	return expressao
}

func GetLista(fila *FilaExpressoes) []arvore.Expressao {
	listaExpressoes := make([]arvore.Expressao,fila.Tamanho)

	for i:= 0;i<fila.Tamanho;i++{
		listaExpressoes[i] = fila.GetItem()
	}

	return listaExpressoes
}

//Caso especifico dos cases do switch
type FilaCases struct {
	Tamanho  int
	primeiro *ItemCaso
	ultimo   *ItemCaso
}

type ItemCaso struct {
	item    *arvore.Case
	proximo *ItemCaso
}

func NewFilaCasos() *FilaCases {
	fila := &FilaCases{Tamanho: 0}
	return fila
}

func (fl *FilaCases) AddItem(expr *arvore.Case) {
	novo := &ItemCaso{item: expr}
	if fl.Tamanho == 0 {
		fl.primeiro = novo
		fl.ultimo = novo
		fl.Tamanho++
		return
	}
	fl.ultimo.proximo = novo
	fl.ultimo = novo
	fl.Tamanho++
}
func (fl *FilaCases) GetItem() *arvore.Case {
	expressao := fl.primeiro.item
	fl.primeiro = fl.primeiro.proximo
	return expressao
}
