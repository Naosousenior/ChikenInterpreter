package objetos

import (
	//arv "TestesInterpreter/parsing/arvore"

	tool "ChikenInterpreter/ferramentas"
	"fmt"
	"strings"
)

type ObjArray struct {
	ArrayList  []ObjetoBase
	Capacidade int
	Tamanho    int
}

func (l *ObjArray) Tipo() TipoObjeto { return LISTA }
func (l *ObjArray) Inspecionar() string {
	if l.Tamanho == 0 {
		return "[]"
	}
	parts := make([]string, l.Tamanho)

	for i := 0; i < l.Tamanho; i++ {
		parts[i] = l.ArrayList[i].Inspecionar()
	}

	parts[0] = "[" + parts[0]
	parts[l.Tamanho-1] += "]"

	return strings.Join(parts, ",")
}

func (l *ObjArray) OpInfixo(op string, valor ObjetoBase) ObjetoBase {
	switch op {
	case "+":
		if valor.Tipo() == LISTA {
			concatenado := valor.(*ObjArray)
			nova_lista := append(l.ArrayList, concatenado.ArrayList[:]...)

			return &ObjArray{ArrayList: nova_lista, Tamanho: len(nova_lista), Capacidade: len(nova_lista)}

		} else {
			nova_lista := append(l.ArrayList, valor)
			return &ObjArray{ArrayList: nova_lista, Tamanho: len(nova_lista), Capacidade: len(nova_lista)}
		}

	case ":":
		for _, obj := range l.ArrayList {
			if obj.OpInfixo("==", valor) == OBJ_TRUE {
				return OBJ_TRUE
			}
		}

		return OBJ_FALSE

	case "[":
		if valor.Tipo() == INTEIRO {
			index := valor.(*ObjInteiro).Valor
			if index >= l.Tamanho {
				return geraErro(fmt.Sprintf("Posicao %d inexistente", index))
			}

			if index < 0 {
				index += l.Tamanho
				if index < 0 {
					return geraErro(fmt.Sprintf("Posicao %d inexistente", index-l.Tamanho))
				}
			}

			return l.ArrayList[index]
		} else {
			return geraErro("Listas so podem ser indexadas por numeros inteiros")
		}

	default:
		return geraErro(fmt.Sprintf("Operacao %s incompativel com listas", op))
	}
}

func (l *ObjArray) OpPrefixo(op string) ObjetoBase {
	return geraErro("Objetos lista sao incompativeis com operadores de prefixo")
}

func (l *ObjArray) GetPropriedade(propri string) ObjetoBase {
	switch propri {
	case "length":
		return &ObjInteiro{Valor: l.Tamanho}

	case "capacity":
		return &ObjInteiro{Valor: l.Capacidade}

	default:
		return geraErro("Propriedade " + propri + " inexistente.")
	}
}
func (l *ObjArray) SetPropriedade(propri string, valor ObjetoBase) ObjetoBase {

	switch propri {
	case "length":
		if valor.Tipo() == INTEIRO {
			novoLimite := valor.(*ObjInteiro).Valor
			if novoLimite == l.Capacidade {
				return OBJ_NONE
			}
			novoArray := make([]ObjetoBase, novoLimite)

			if novoLimite > l.Capacidade {
				copy(novoArray[:], l.ArrayList[:])

				i := l.Tamanho
				for i < novoLimite {
					novoArray[i] = OBJ_NONE
					i++
				}
			} else {
				copy(novoArray[:], l.ArrayList[:novoLimite])
				l.Tamanho = novoLimite
			}

			l.ArrayList = novoArray
			l.Capacidade = novoLimite
			return OBJ_NONE
		} else {
			return geraErro("O tamanho da lista deve ser representado por um INTEIRO.")
		}

	default:
		return geraErro("Propriedade nao encontrada, ou não pode ser alterada")
	}
}

func (l *ObjArray) SetIndex(index ObjetoBase, valor ObjetoBase) ObjetoBase {
	if index.Tipo() == INTEIRO {
		posicao := index.(*ObjInteiro).Valor

		if posicao+1 > l.Capacidade {
			return geraErro(fmt.Sprintf("Posicao %d inexistente.", posicao))
		}

		if posicao < 0 {
			posicao = l.Tamanho + posicao

			if posicao < 0 {
				return geraErro(fmt.Sprintf("Posicao %d inexistente", posicao-l.Tamanho))
			}
		}

		if posicao > l.Tamanho {
			l.Tamanho = posicao + 1
		}

		l.ArrayList[posicao] = valor

		return OBJ_NONE
	} else {
		return geraErro("Listas sao enumeradas, o indexador precisa ser um inteiro.")
	}
}
func (l *ObjArray) Iterar(pos int) ObjetoBase {
	if pos < l.Tamanho {
		return l.ArrayList[pos]
	}

	return OBJ_BREAK
}

type ObjDict struct {
	Dict   map[string]ObjetoBase
	Chaves []string
}

func (obj *ObjDict) Tipo() TipoObjeto { return DICT }
func (obj *ObjDict) Inspecionar() string {

	if len(obj.Dict) < 1 {
		return "Dict: {}"
	}

	parts := make([]string, len(obj.Dict)+2)

	parts[0] = "Dict: {"

	i := 0
	for hash := range obj.Dict {
		parts[i+1] = hash + ","
		i++
	}

	parts[len(parts)-1] = "}"

	return strings.Join(parts, "")
}

func (obj *ObjDict) OpPrefixo(op string) ObjetoBase {
	return geraErro("Dicts não suportam o operador " + op)
}

func (obj *ObjDict) OpInfixo(op string, dir ObjetoBase) ObjetoBase {
	switch op {
	case "+":
		if dict, ok := dir.(*ObjDict); ok {
			mapa := dict.Dict

			novoMapa := make(map[string]ObjetoBase)

			i := 0

			for index, valor := range obj.Dict {
				novoMapa[index] = valor
				i++
			}

			for index, valor := range mapa {
				novoMapa[index] = valor
				i++
			}

			i = 0
			chaves := make([]string, len(novoMapa))
			for chave := range novoMapa {
				chaves[i] = chave
				i++
			}

			return &ObjDict{Dict: novoMapa, Chaves: chaves}
		}

		return geraErro("DICT suporta apenas adição com outra DICT")

	case "[":
		hash := string(dir.Tipo()) + ": " + dir.Inspecionar()

		res, ok := obj.Dict[hash]

		if !ok {
			return geraErro(fmt.Sprintf("Chave %s inexiste no objeto %s", hash, obj.Inspecionar()))
		}

		return res

	default:
		return geraErro(fmt.Sprintf("Operação %s não suportada por DICT", op))
	}
}

func (obj *ObjDict) GetPropriedade(propriedade string) ObjetoBase {
	switch propriedade {
	case "length":
		return &ObjInteiro{Valor: len(obj.Dict)}
	case "keys":
		keys := make([]ObjetoBase, len(obj.Dict))

		i := 0
		for key := range obj.Dict {
			keys[i] = &ObjTexto{Valor: key}
			i++
		}

		return &ObjArray{ArrayList: keys, Tamanho: len(keys), Capacidade: len(keys)}

	default:
		return geraErro(fmt.Sprintf("Atributo %s não encontrado", propriedade))
	}
}

func (obj *ObjDict) SetPropriedade(propriedade string, valor ObjetoBase) ObjetoBase {
	return geraErro(fmt.Sprintf("Propriedade %s apenas leitura ou inexistente.", propriedade))
}

func (obj *ObjDict) SetIndex(index ObjetoBase, valor ObjetoBase) ObjetoBase {
	hash := string(index.Tipo()) + ": " + index.Inspecionar()
	if _, ok := obj.Dict[hash]; !ok {
		obj.Chaves = append(obj.Chaves, hash)
	}

	obj.Dict[hash] = valor

	return OBJ_NONE
}
func (obj *ObjDict) Iterar(index int) ObjetoBase {
	if index >= len(obj.Chaves) {
		return OBJ_BREAK
	}

	return obj.Dict[obj.Chaves[index]]
}

type Classe struct {
	Supers          []*Classe
	Construtor      *Metodo
	ObjModel        *ObjetoUser
	AtribbProtegido tool.Conjunto
	AtributosClass  Propriedades
}

type Metodo struct {
	Classe     *Classe
	Objeto     *ObjetoUser
	Funcao     *ObjFuncao
	Parametros []string
}

func (cl *Classe) Tipo() TipoObjeto { return CLASSE }
func (cl *Classe) Inspecionar() string {
	
	//fmt.Println(cl.ObjModel.Publicas["nome"].Inspecionar())

	parts := make([]string, len(cl.AtributosClass))

	i := 0
	for text := range cl.AtributosClass {
		parts[i] = text
		i++
	}

	return fmt.Sprintf("Class: {%s}", strings.Join(parts, ","))
}

func (cl *Classe) OpInfixo(op string, dir ObjetoBase) ObjetoBase {
	return geraErro("Classes nao suportam operacoes")
}
func (cl *Classe) OpPrefixo(op string) ObjetoBase {
	return geraErro("Classes nao suportam operacoes")
}
func (cl *Classe) GetPropriedade(propri string) ObjetoBase {
	var res ObjetoBase
	res, ok := cl.AtributosClass[propri]
	if !ok {

		return geraErro(fmt.Sprintf("Propriedade %s não encontrada", propri))
	}

	return res
}
func (cl *Classe) SetPropriedade(propri string, valor ObjetoBase) ObjetoBase {
	if _, ok := cl.AtributosClass[propri]; !ok {
		return geraErro(fmt.Sprintf("Propriedade %s não encontrada", propri))
	}

	cl.AtributosClass[propri] = valor

	return OBJ_NONE
}

func (m *Metodo) Tipo() TipoObjeto { return FUNCAO_OBJ }
func (f *Metodo) Inspecionar() string {
	partes := make([]string, len(f.Parametros)+2)

	partes[0] = "method ("

	for i, v := range f.Parametros {
		partes[i+1] = v + ", "
	}

	partes[len(partes)-1] = ")"

	return strings.Join(partes, "")
}

func (f *Metodo) OpInfixo(op string, dir ObjetoBase) ObjetoBase {
	return geraErro("Funcoes nao suportam operacoes")
}

func (f *Metodo) OpPrefixo(op string) ObjetoBase {
	return geraErro("Funcoes nao suportam operacoes")
}
func (f *Metodo) GetPropriedade(propri string) ObjetoBase {
	return geraErro("Funcoes nao possuem propriedades")
}
func (f *Metodo) SetPropriedade(propri string, valor ObjetoBase) ObjetoBase {
	return geraErro("Funcoes nao possuem propriedades")
}

type ObjetoUser struct {
	ClasseMae  *Classe
	Publicas   Propriedades
	Protegidos Propriedades
	Privadas   map[*Classe]Propriedades
}

func (obj *ObjetoUser) Tipo() TipoObjeto { return OBJETO }
func (obj *ObjetoUser) Inspecionar() string {
	return fmt.Sprintf("Object class = %s, address = %d", obj.ClasseMae.Inspecionar(), &obj)
}
func (obj *ObjetoUser) OpPrefixo(op string) ObjetoBase {
	return &ObjTexto{Valor: "Nao implementado"}
}
func (obj *ObjetoUser) OpInfixo(op string, dir ObjetoBase) ObjetoBase {
	return &ObjTexto{Valor: "Nao implementado"}
}
func (obj *ObjetoUser) GetPropriedade(propriedade string) ObjetoBase {
	res, ok := obj.Publicas[propriedade]

	if !ok {
		return geraErro("Propriedade " + propriedade + " nao encontrada")
	}

	if metodo, ok := res.(*Metodo); ok {
		metodo.Objeto = obj

		return metodo
	}

	return res
}

func (obj *ObjetoUser) SetPropriedade(prorpri string, valor ObjetoBase) ObjetoBase {
	if _, ok := obj.Publicas[prorpri]; ok {
		obj.Publicas[prorpri] = valor

		return OBJ_NONE
	}

	return geraErro(fmt.Sprintf("Propriedade %s nao encontrada", prorpri))
}

func (obj *ObjetoUser) Get(propriedade string, ambiente *Ambiente) ObjetoBase {
	if ambiente.Classe.AtribbProtegido.Tem(propriedade) {
		res := obj.Protegidos[propriedade]
		if metodo, ok := res.(*Metodo); ok {
			metodo.Objeto = obj
			return metodo
		}

		return res
	} else if classe, ok := obj.Privadas[ambiente.Classe]; ok {
		res, ok := classe[propriedade]

		if !ok {
			return obj.GetPropriedade(propriedade)
		}

		if metodo, ok := res.(*Metodo); ok {
			metodo.Objeto = obj
			return metodo
		}

		return res
	}

	return obj.GetPropriedade(propriedade)
}

func (obj *ObjetoUser) Set(propriedade string, valor ObjetoBase, ambiente *Ambiente) ObjetoBase {
	if _, ok := obj.Protegidos[propriedade]; ok {
		obj.Protegidos[propriedade] = valor
		return OBJ_NONE
	} else if classe, ok := obj.Privadas[ambiente.Classe]; ok {
		_, ok := classe[propriedade]

		if !ok {
			return obj.SetPropriedade(propriedade, valor)
		}

		classe[propriedade] = valor
		return OBJ_NONE
	}

	return obj.SetPropriedade(propriedade, valor)
}
