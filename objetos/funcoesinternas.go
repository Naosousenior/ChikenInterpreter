package objetos

import (
	"bufio"
	"fmt"
	"os"
)

type (
	Funcao func([]ObjetoBase) ObjetoBase
)
type FuncaoInterna struct {
	nome   string
	Funcao Funcao
}

func (f *FuncaoInterna) Tipo() TipoObjeto { return FUNCAO_INTERNA }
func (f *FuncaoInterna) Inspecionar() string {
	return fmt.Sprintf("Funcao interna '%s'", f.nome)
}
func (f *FuncaoInterna) OpInfixo(op string, dir ObjetoBase) ObjetoBase {
	return geraErro("Funcoes internas nao suportam operacoes")
}
func (f *FuncaoInterna) OpPrefixo(op string) ObjetoBase {
	return geraErro("Funcoes internas nao suportam operacoes")
}
func (f *FuncaoInterna) GetPropriedade(propri string) ObjetoBase {
	return geraErro("Funcoes internas nao possuem propriedades")
}
func (f *FuncaoInterna) SetPropriedade(propri string, valor ObjetoBase) ObjetoBase {
	return geraErro("Funcoes internas nao possuem propriedades")
}

func NewFuncaoInterna(funcao Funcao, nome string) *FuncaoInterna {
	return &FuncaoInterna{Funcao: funcao, nome: nome}
}

func Write(args []ObjetoBase) ObjetoBase {
	for _, obj := range args {
		fmt.Print(obj.Inspecionar())
	}
	return OBJ_NONE
}

var leitor = bufio.NewScanner(os.Stdin)

func Read(args []ObjetoBase) ObjetoBase {
	Write(args)

	if leitor.Scan() {
		return &ObjTexto{Valor: leitor.Text()}
	}

	return OBJ_NONE
}
