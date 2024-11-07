package evaluation

import (
	tool "ChikenInterpreter/ferramentas"
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

func chamaMetodo(funcao *obj.ObjFuncao, ambiente *obj.Ambiente) obj.ObjetoBase {
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
		novoAmb := passaParametros(args, objeto.Funcao)

		novoAmb.Objeto = objeto.Objeto
		novoAmb.Classe = objeto.Classe

		return chamaMetodo(objeto.Funcao, novoAmb)

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

func avaliaObject(expressao *arv.ExpressaoObjeto, classe *obj.Classe) (*obj.ObjetoUser, obj.ObjetoBase) {
	var valor obj.ObjetoBase
	ambiente := obj.NewAmbiente()
	ambiente.Classe = classe
	novoObjeto := &obj.ObjetoUser{
		ClasseMae:  ambiente.Classe,
		Publicas:   make(obj.Propriedades),
		Protegidos: make(obj.Propriedades),
		Privadas:   make(map[*obj.Classe]obj.Propriedades),
	}

	constroiObjeto(novoObjeto, CLASSMAE)

	for propriedade, exprV := range expressao.Atributos {
		valor = Avaliar(exprV, ambiente)

		if valor.Tipo() == obj.ERRO {
			return nil, valor
		}

		if valor.Tipo() == obj.FUNCAO_OBJ {
			funcao, _ := valor.(*obj.ObjFuncao)
			params := make([]string, len(funcao.Parametros))

			for i, text := range funcao.Parametros {
				params[i] = text.Nome
			}

			valor = &obj.Metodo{Classe: ambiente.Classe, Funcao: funcao, Parametros: params}
		}

		novoObjeto.Publicas[propriedade] = valor
	}

	return novoObjeto, nil
}

func constroiObjeto(objeto *obj.ObjetoUser, classeMae *obj.Classe) {
	if classeMae == CLASSMAE {
		for propriedade, valor := range CLASSMAE.ObjModel.Publicas {
			objeto.Publicas[propriedade] = valor
		}

		return
	} else {
		for i := len(classeMae.Supers) - 1; i >= 0; i++ {
			constroiObjeto(objeto, classeMae.Supers[i])
		}

		passaAtributos(&objeto.Publicas, &classeMae.ObjModel.Publicas)
		passaAtributos(&objeto.Protegidos, &classeMae.ObjModel.Protegidos)

		for nome, valor := range classeMae.ObjModel.Privadas[classeMae] {
			objeto.Privadas[classeMae][nome] = valor
		}
	}
}

//func chamaConstrutor(objeto *obj.ObjetoUser,classe *obj.Classe) (*obj.ObjetoUser,*obj.ObjErro)

func passaAtributos(atribRec *obj.Propriedades, atribMen *obj.Propriedades) {
	for name, value := range *atribMen {
		(*atribRec)[name] = value
	}
}

func avaliaClasse(noClass *arv.ExpressaoClass, superClasses []*obj.Classe, ambiente *obj.Ambiente) obj.ObjetoBase {
	novaClasse := &obj.Classe{Supers: superClasses}

	//prepara o modelo
	modeloObjeto := &obj.ObjetoUser{
		ClasseMae:  novaClasse,
		Publicas:   make(obj.Propriedades),
		Protegidos: make(obj.Propriedades),
		Privadas:   make(map[*obj.Classe]obj.Propriedades),
	}

	for i := len(superClasses) - 1; i >= 0; i-- {
		construirModeloObjeto(modeloObjeto, superClasses[i])
	}

	erro, construtor := addAtributosAtuais(modeloObjeto, noClass, ambiente)

	if erro != nil {
		return erro
	}

	//com o modelo pronto, basta adicionar ele no objeto classe
	novaClasse.ObjModel = modeloObjeto

	//agora, vamos adicionar um objeto construtor
	novaClasse.Construtor = construtor

	//por fim, vamos adicionar os atributos de classe
	novaClasse.AtributosClass, erro = getAtributosClasse(noClass, superClasses, ambiente)

	if erro != nil {
		return erro
	}

	novaClasse.AtribbProtegido = getAtribbProtegidos(novaClasse)

	return novaClasse
}

func construirModeloObjeto(modelo *obj.ObjetoUser, classeModelo *obj.Classe) {

	for chave, valor := range classeModelo.ObjModel.Publicas {
		modelo.Publicas[chave] = valor
	}

	for chave, valor := range classeModelo.ObjModel.Protegidos {
		modelo.Protegidos[chave] = valor
	}

	for chave, valor := range classeModelo.ObjModel.Privadas[classeModelo] {
		modelo.Privadas[classeModelo][chave] = valor
	}
}

func addAtributosAtuais(modelo *obj.ObjetoUser, expreClasse *arv.ExpressaoClass, ambiente *obj.Ambiente) (obj.ObjetoBase, *obj.Metodo) {
	var construtor *obj.Metodo = nil

	for chave, valor := range expreClasse.AtribPub {
		atributo := Avaliar(valor, ambiente)

		if metodo, ok := atributo.(*obj.ObjFuncao); ok {
			//o nome de atributo "new_object"
			//identifica um construtor
			//obviamente isso so é válido
			//se for um atributo público
			if chave == "new_object" {
				construtor = newMetodo(modelo.ClasseMae, metodo)
				continue
			}

			atributo = newMetodo(modelo.ClasseMae, metodo)

		} else if atributo.Tipo() == obj.ERRO {
			return atributo, nil
		}

		modelo.Publicas[chave] = atributo
	}

	for chave, valor := range expreClasse.AtribPro {
		atributo := Avaliar(valor, ambiente)

		if metodo, ok := atributo.(*obj.ObjFuncao); ok {
			atributo = newMetodo(modelo.ClasseMae, metodo)
		} else if atributo.Tipo() == obj.ERRO {
			return atributo, nil
		}

		modelo.Protegidos[chave] = atributo
	}

	for chave, valor := range expreClasse.AtribPriv {
		atributo := Avaliar(valor, ambiente)

		if metodo, ok := atributo.(*obj.ObjFuncao); ok {
			atributo = newMetodo(modelo.ClasseMae, metodo)
		} else if atributo.Tipo() == obj.ERRO {
			return atributo, nil
		}

		modelo.Privadas[ambiente.Classe][chave] = atributo
	}

	return nil, construtor
}

func getAtribbProtegidos(classe *obj.Classe) tool.Conjunto {

	protegidos := make(tool.Conjunto)

	for _, super := range classe.Supers {
		protegidos.Copiar(super.AtribbProtegido)
	}

	for chave := range classe.ObjModel.Protegidos {
		protegidos.Add(chave)
	}

	return protegidos
}

func newMetodo(classe *obj.Classe, funcao *obj.ObjFuncao) *obj.Metodo {
	novoMetodo := &obj.Metodo{Classe: classe, Funcao: funcao, Parametros: make([]string, len(funcao.Parametros))}

	for i, atrib := range funcao.Parametros {
		novoMetodo.Parametros[i] = atrib.Nome
	}

	return novoMetodo
}

func getAtributosClasse(expreClass *arv.ExpressaoClass, superClasses []*obj.Classe, amb *obj.Ambiente) (obj.Propriedades, obj.ObjetoBase) {
	atribClasse := make(obj.Propriedades)

	for _, super := range superClasses {
		for chave, valor := range super.AtributosClass {
			atribClasse[chave] = valor
		}
	}

	for chave, expre := range expreClass.AtribClass {
		atributo := Avaliar(expre, amb)

		if atributo.Tipo() == obj.ERRO {
			return nil, atributo
		}

		atribClasse[chave] = atributo
	}

	return atribClasse, nil
}
