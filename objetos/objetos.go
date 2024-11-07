package objetos

type TipoObjeto string
type Propriedades map[string]ObjetoBase

const (
	BFF = "BFF"

	INTEIRO        = "INT"
	REAL           = "REAL"
	TEXTO          = "TEXTO"
	BOOLEANO       = "BOOL"
	NONE           = "NONE"
	VALOR_RETORNO  = "RETORNO"
	INSTRUCAO      = "INSTRUCAO"
	FUNCAO_OBJ     = "FUNCAO"
	FUNCAO_INTERNA = "FUNC INTERNA"
	ERRO           = "ERRO"

	LISTA  = "LISTA"
	DICT   = "DICT"
	CLASSE = "CLASSE"
	OBJETO = "OBJETO"
)

var (
	OBJ_TRUE     = &ObjBool{Valor: true}
	OBJ_FALSE    = &ObjBool{Valor: false}
	OBJ_NONE     = &ObjNone{}
	OBJ_BREAK    = &ObjInstrucao{Instru: "BREAK"}
	OBJ_CONTINUE = &ObjInstrucao{Instru: "CONTINUE"}
)

type ObjetoBase interface {
	Tipo() TipoObjeto
	Inspecionar() string
	OpInfixo(string, ObjetoBase) ObjetoBase
	OpPrefixo(string) ObjetoBase
	GetPropriedade(string) ObjetoBase
	SetPropriedade(string, ObjetoBase) ObjetoBase
}

type ObjetoIndexavel interface {
	ObjetoBase
	SetIndex(ObjetoBase, ObjetoBase) ObjetoBase
	Iterar(int) ObjetoBase
}
