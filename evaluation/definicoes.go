package evaluation

import (
	obj "ChikenInterpreter/objetos"
	"fmt"
)

var (
	CLASSMAE = NewClasseMae()
)

func NewClasseMae() *obj.Classe {
	var classe *obj.Classe

	propriPadrao := map[string] obj.ObjetoBase {
		"__type__": &obj.ObjTexto{Valor: string(obj.OBJETO)},
	}

	objeto := &obj.ObjetoUser{ClasseMae: classe,Publicas: propriPadrao}
	classe = &obj.Classe{ObjModel: objeto}

	return classe
}

type Bff struct {
	FuncoesInternas map[string] *obj.FuncaoInterna
}

func (bff *Bff) Tipo() obj.TipoObjeto { return obj.BFF }
func (bff *Bff) Inspecionar() string {
	return "Objeto Best Friend Forever"
}

func (bff *Bff) OpInfixo(op string, dir obj.ObjetoBase) obj.ObjetoBase {
	return obj.OBJ_NONE
}

func (bff *Bff) OpPrefixo(op string) obj.ObjetoBase {
	return obj.OBJ_NONE
}
func (bff *Bff) GetPropriedade(propri string) obj.ObjetoBase {
	res, ok := bff.FuncoesInternas[propri]

	if !ok {
		return geraErro(fmt.Sprintf("Funcao interna %s desconhecida",propri))
	}

	return res
}
func (bff *Bff) SetPropriedade(propri string, valor obj.ObjetoBase) obj.ObjetoBase {
	return geraErro("Você não deve e nem pode modificar as propriedades do seu BFF.")
}

func NewBff() *Bff {
	bff := &Bff{}
	bff.FuncoesInternas = make(map[string]*obj.FuncaoInterna)

	bff.registraFuns()

	return bff
}

func (bff *Bff) registraFuns() {
	bff.FuncoesInternas["write"] = obj.NewFuncaoInterna(obj.Write,"write")
	bff.FuncoesInternas["read"] = obj.NewFuncaoInterna(obj.Read,"read")
}