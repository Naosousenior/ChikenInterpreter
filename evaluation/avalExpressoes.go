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

	if resultado.Tipo == obj.RETURN || resultado.Tipo == obj.ERROR {

		return resultado.Resultado
	} else {
		return obj.OBJ_NONE
	}
}

func chamaMetodo(metodo *obj.Metodo, argumentos []obj.ObjetoBase) obj.ObjetoBase {
	novoAmb := passaParametros(argumentos, metodo.Funcao)

	novoAmb.Objeto = metodo.Objeto
	novoAmb.Classe = metodo.Classe
	resultado := avaliaInstrucoes(metodo.Funcao.BlocoInstrucoes.Instrucoes, novoAmb)

	if resultado.Tipo == obj.RETURN || resultado.Tipo == obj.ERROR {

		return resultado.Resultado
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
		if args[0].Tipo() == obj.EXCECAO {
			return args[0]
		}
	}

	switch objeto := objeto.(type) {
	case *obj.ObjFuncao:

		return chamaFuncao(objeto, args)

	case *obj.FuncaoInterna:

		return objeto.Funcao(args)

	case *obj.Metodo:
		//preparacao do ambiente
		

		return chamaMetodo(objeto, args)

	case *obj.Classe:
		//implementacao de uma instanciacao de objeto


		//primeiro, verificamos se o construtor da classe foi definido pelo usuário
		novo_objeto := instanciarNovoObjeto(objeto)
		if objeto.Construtor != nil {
			construtor := objeto.Construtor
			construtor.Objeto = novo_objeto
			
			resultado := chamaMetodo(construtor,args)

			if resultado.Tipo() == obj.EXCECAO {
				return resultado
			}
		}

		return novo_objeto
		
	default:
		return geraErro(fmt.Sprintf("O objeto %s não pode ser chamado", objeto.Inspecionar()))
	}
}

func instanciarNovoObjeto(classeMae *obj.Classe) *obj.ObjetoUser {
	novoObj := obj.ObjetoUser{
		ClasseMae: classeMae,
		Publicas: make(obj.Propriedades),
		Protegidos: make(obj.Propriedades),
		Privadas: make(map[*obj.Classe]obj.Propriedades),
	}

	passaAtributos(&novoObj,classeMae.ObjModel)

	return &novoObj
}

func passaAtributos(receptor *obj.ObjetoUser, modelo *obj.ObjetoUser) {
	for name, value := range modelo.Publicas {
		receptor.Publicas[name] = value
	}

	for name, value := range modelo.Protegidos {
		receptor.Protegidos[name] = value
	}

	classeMae := modelo.ClasseMae

	for _,super := range classeMae.Supers {
		receptor.Privadas[super] = make(obj.Propriedades)
		for nome,atrib := range modelo.Privadas[super]{
			receptor.Privadas[super][nome] = atrib
		}
	}

	for nome,atrib := range modelo.Privadas[classeMae]{
		receptor.Privadas[classeMae][nome] = atrib
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
	i := 0

	for i < len(expressao.Chaves) {
		key = AvaliaExpressao(expressao.Chaves[i], ambiente)
		value = AvaliaExpressao(expressao.Valores[i], ambiente)

		if key.Tipo() == obj.EXCECAO {
			return key
		}

		if value.Tipo() == obj.EXCECAO {
			return value
		}

		hash = FalsoHash(key)

		dict[hash] = value

		i++
	}

	return &obj.ObjDict{Dict: dict}
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

	for propriedade, valor := range CLASSMAE.ObjModel.Publicas {
		novoObjeto.Publicas[propriedade] = valor
	}

	for propriedade, exprV := range expressao.Atributos {
		valor = AvaliaExpressao(exprV, ambiente)

		if valor.Tipo() == obj.EXCECAO {
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

/*

Metodo ultrapassado

func constroiObjeto(objeto *obj.ObjetoUser, classeMae *obj.Classe) {
	if classeMae == CLASSMAE {
		for propriedade, valor := range CLASSMAE.ObjModel.Publicas {
			objeto.Publicas[propriedade] = valor
		}
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

*/

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

	modelo.Privadas[classeModelo] = make(obj.Propriedades)
	for chave, valor := range classeModelo.ObjModel.Privadas[classeModelo] {
		modelo.Privadas[classeModelo][chave] = valor
	}
}

func addAtributosAtuais(modelo *obj.ObjetoUser, expreClasse *arv.ExpressaoClass, ambiente *obj.Ambiente) (obj.ObjetoBase, *obj.Metodo) {
	var construtor *obj.Metodo = nil

	for chave, valor := range expreClasse.AtribPub {
		atributo := AvaliaExpressao(valor, ambiente)

		if metodo, ok := atributo.(*obj.ObjFuncao); ok {
			//o nome de atributo "new_object"
			//identifica um construtor
			//obviamente isso so é válido
			//se for um atributo público
			if chave == "new_object" {
				fmt.Println("Achei o problema")
				construtor = newMetodo(modelo.ClasseMae, metodo)
				continue
			}

			atributo = newMetodo(modelo.ClasseMae, metodo)

		} else if atributo.Tipo() == obj.EXCECAO {
			return atributo, nil
		}

		modelo.Publicas[chave] = atributo
	}

	for chave, valor := range expreClasse.AtribPro {
		atributo := AvaliaExpressao(valor, ambiente)

		if metodo, ok := atributo.(*obj.ObjFuncao); ok {
			atributo = newMetodo(modelo.ClasseMae, metodo)
		} else if atributo.Tipo() == obj.EXCECAO {
			return atributo, nil
		}

		modelo.Protegidos[chave] = atributo
	}

	modelo.Privadas[modelo.ClasseMae] = make(obj.Propriedades)
	for chave, valor := range expreClasse.AtribPriv {
		atributo := AvaliaExpressao(valor, ambiente)

		if metodo, ok := atributo.(*obj.ObjFuncao); ok {
			atributo = newMetodo(modelo.ClasseMae, metodo)
		} else if atributo.Tipo() == obj.EXCECAO {
			return atributo, nil
		}

		fmt.Printf("chave: %s, classe: %s, atributo: %s",chave,modelo.ClasseMae.Inspecionar(),atributo.Inspecionar())
		
		modelo.Privadas[modelo.ClasseMae][chave] = atributo
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
		atributo := AvaliaExpressao(expre, amb)

		if atributo.Tipo() == obj.EXCECAO {
			return nil, atributo
		}

		atribClasse[chave] = atributo
	}

	return atribClasse, nil
}
