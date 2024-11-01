package evaluation

import (
	obj "ChikenInterpreter/objetos"
	arv "ChikenInterpreter/parsing/arvore"

	"fmt"
)

func chamaFuncao(funcao *obj.ObjFuncao, argumentos []obj.ObjetoBase) obj.ObjetoBase {
	novoAmb := passaParametros(argumentos, funcao)

	resultado := avaliaInstrucoes(funcao.BlocoInstrucoes.Instrucoes, novoAmb)

	if retorno, ok := resultado.(*obj.ObjReturn); ok {

		return retorno.Valor
	} else if resultado.Tipo() == obj.ERRO {
		return resultado
	} else {
		return obj.OBJ_NONE
	}
}

func chamaMetodo(funcao *obj.ObjFuncao,ambiente *obj.Ambiente) obj.ObjetoBase{
	//fmt.Println(ambiente.Objeto.Protegidos)
	resultado := avaliaInstrucoes(funcao.BlocoInstrucoes.Instrucoes, ambiente)

	if retorno, ok := resultado.(*obj.ObjReturn); ok {

		return retorno.Valor
	} else if resultado.Tipo() == obj.ERRO {
		return resultado
	} else {
		return obj.OBJ_NONE
	}
}

func passaParametros(args []obj.ObjetoBase, funcao *obj.ObjFuncao) *obj.Ambiente {
	novoAmb := obj.NewAmbienteInterno(funcao.Amb)

	for i, id := range funcao.Parametros {
		novoAmb.AddArgs(id.Nome, args[i])
	}

	return novoAmb
}

func avaliaChamada(noCall *arv.CallFun, objeto obj.ObjetoBase, ambiente *obj.Ambiente) obj.ObjetoBase {
	var args []obj.ObjetoBase
	if noCall.Argumentos != nil {
		args = avaliaExpressoes(noCall.Argumentos, ambiente)
	} else {
		args = make([]obj.ObjetoBase, 0)
	}

	if len(args) > 0 {
		if args[0].Tipo() == obj.ERRO {
			return args[0]
		}
	}

	switch objeto := objeto.(type) {
	case *obj.ObjFuncao:

		return chamaFuncao(objeto, args)

	case *obj.FuncaoInterna:

		return objeto.Funcao(args)

	case *obj.Metodo:
		novoAmb := passaParametros(args,objeto.Funcao)

		novoAmb.Objeto = objeto.Objeto
		novoAmb.Classe = objeto.Classe

		return chamaMetodo(objeto.Funcao,novoAmb) 

	default:
		return geraErro(fmt.Sprintf("O objeto %s não pode ser chamado", objeto.Inspecionar()))
	}
}

func avaliaIdentificador(nome string, ambiente *obj.Ambiente) obj.ObjetoBase {
	res, err := ambiente.GetVar(nome)
	if err {
		return res
	} else {
		return geraErro(fmt.Sprintf("Identificador %s nao declarado no contexto. ", nome))
	}
}

func avaliaInfixo(operador string,
	vEsq, vDir obj.ObjetoBase) obj.ObjetoBase {
	if operador == "==" || operador == "!=" {
		if vEsq.Tipo() == obj.REAL || vEsq.Tipo() == obj.INTEIRO && vDir.Tipo() == obj.REAL || vDir.Tipo() == obj.INTEIRO {
			return vEsq.OpInfixo(operador, vDir)
		}

		if operador == "==" {
			if vEsq.Tipo() != vDir.Tipo() {
				return obj.OBJ_FALSE
			}
		} else {
			if vEsq.Tipo() == vDir.Tipo() {
				return obj.OBJ_FALSE
			}
		}
	}

	return vEsq.OpInfixo(operador, vDir)
}

func avaliaDict(expressao *arv.ExpressaoDict, ambiente *obj.Ambiente) obj.ObjetoBase {
	var key obj.ObjetoBase
	var value obj.ObjetoBase
	var hash string

	dict := make(map[string]obj.ObjetoBase)
	chaves := make([]string, len(expressao.Chaves))
	i := 0

	for i < len(expressao.Chaves) {
		key = Avaliar(expressao.Chaves[i], ambiente)
		value = Avaliar(expressao.Valores[i], ambiente)

		if key.Tipo() == obj.ERRO {
			return key
		}

		if value.Tipo() == obj.ERRO {
			return value
		}

		hash = FalsoHash(key)

		dict[hash] = value
		chaves[i] = hash

		i++
	}

	return &obj.ObjDict{Dict: dict, Chaves: chaves}
}

func avaliaObject(expressao *arv.ExpressaoObjeto, classe *obj.Classe) (*obj.ObjetoUser,obj.ObjetoBase) {
	var valor obj.ObjetoBase
	ambiente := obj.NewAmbiente()
	ambiente.Classe = classe
	novoObjeto := &obj.ObjetoUser{
		ClasseMae: ambiente.Classe,
		Publicas: make(obj.Propriedade),
		Protegidos: make(obj.Propriedade),
		Privadas: make(map[*obj.Classe]obj.Propriedade),
	}

	constroiObjeto(novoObjeto,CLASSMAE)

	for propriedade, exprV := range expressao.Atributos {
		valor = Avaliar(exprV,ambiente)

		if valor.Tipo() == obj.ERRO {
			return nil, valor
		}

		if valor.Tipo() == obj.FUNCAO_OBJ {
			funcao,_ := valor.(*obj.ObjFuncao)
			params := make([]string,len(funcao.Parametros))

			for i,text := range funcao.Parametros{
				params[i] = text.Nome
			}

			valor = &obj.Metodo{Classe: ambiente.Classe,Funcao: funcao,Parametros: params}
		}

		novoObjeto.Publicas[propriedade] = valor
	}

	return novoObjeto,nil
}

func constroiObjeto(objeto *obj.ObjetoUser,classeMae *obj.Classe) {
	if classeMae == CLASSMAE {
		for propriedade, valor := range CLASSMAE.ObjModel.Publicas {
			objeto.Publicas[propriedade] = valor
		}

		return
	} else {
		for i:= len(classeMae.Supers)-1;i >= 0; i++ {
			constroiObjeto(objeto,classeMae.Supers[i])
		}

		passaAtributos(&objeto.Publicas,&classeMae.ObjModel.Publicas)
		passaAtributos(&objeto.Protegidos,&classeMae.ObjModel.Protegidos)
		
		for nome,valor := range classeMae.ObjModel.Privadas[classeMae] {
			objeto.Privadas[classeMae][nome] = valor
		}
	}
}

//func chamaConstrutor(objeto *obj.ObjetoUser,classe *obj.Classe) (*obj.ObjetoUser,*obj.ObjErro)

func passaAtributos(atribRec *obj.Propriedade,atribMen *obj.Propriedade) {
	for name,value := range *atribMen {
		(*atribRec)[name] = value
	}
}

func avaliaClasse(noClass *arv.ExpressaoClass, ambiente *obj.Ambiente) obj.ObjetoBase {
	var temp obj.ObjetoBase
	classe := &obj.Classe{}
	supers,err := passaSuperClasses(noClass.SuperClasses, ambiente)

	if err != nil {
		return err
	}
	classe.Supers = supers

	objetoModel,err := avaliaObject(noClass.Objeto,CLASSMAE)

	if err != nil {
		return err
	}
	classe.ObjModel = objetoModel

	for name,expr := range noClass.AtributosObj {
		temp = Avaliar(expr,ambiente)

		if temp.Tipo() == obj.ERRO {
			return temp
		}

		if funcao,ok := temp.(*obj.ObjFuncao); ok {
			params := make([]string,len(funcao.Parametros))

			for i,text := range funcao.Parametros{
				params[i] = text.Nome
			}

			if name == "new_object" {
				classe.Construtor = &obj.Metodo{Classe: classe,Funcao: funcao,Parametros: params}
				continue
			}

			objetoModel.Publicas[name] = &obj.Metodo{Classe: classe,Funcao: funcao,Parametros: params}

			continue
		}

		objetoModel.Publicas[name] = temp
	}

	classe.AtributosClass = make(obj.Propriedade,len(noClass.AtributosClass))

	for name,expr := range noClass.AtributosClass {
		temp = Avaliar(expr,ambiente)
		if temp.Tipo() == obj.ERRO {
			return temp
		}
		classe.AtributosClass[name] = temp
	}

	return classe
}

func passaSuperClasses(supers []arv.Expressao,amb *obj.Ambiente) ([]*obj.Classe, obj.ObjetoBase) {
	var classeTemp obj.ObjetoBase
	classes := make([]*obj.Classe,len(supers))

	for i,objeto := range supers {
		classeTemp = Avaliar(objeto,amb)

		if classeTemp.Tipo() == obj.ERRO {
			return nil,classeTemp
		}

		if classe,ok := classeTemp.(*obj.Classe); ok {
			classes[i] = classe
			continue
		}

		return nil,geraErro(fmt.Sprintf("Objeto %s não pode ser um super classe",classeTemp.Inspecionar()))
	}

	return classes,nil
}
