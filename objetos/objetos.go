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
	EXCECAO        = "EXCECAO"
	FUNCAO_OBJ     = "FUNCAO"
	FUNCAO_INTERNA = "FUNC INTERNA"

	LISTA  = "LISTA"
	DICT   = "DICT"
	CLASSE = "CLASSE"
	OBJETO = "OBJETO"
)

var (
	OBJ_TRUE  = &ObjBool{Valor: true}
	OBJ_FALSE = &ObjBool{Valor: false}
	OBJ_NONE  = &ObjNone{}
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
	Iterar() chan ObjetoBase
}

type TipoStatus int

const (
	BREAK = iota
	CONTINUE
	RETURN
	ERROR
	TENTATIVA
	DECLARACAO
	DEFINICAO
	ATRIBUICAO
	ITERACAO
	SWITCH
	EXE_INSTRUCAO
	SUPERCALL
	EXPRESSAO
)

var (
	BREAK_ST    = &Status{Tipo: BREAK, Resultado: OBJ_NONE}
	CONTINUE_ST = &Status{Tipo: CONTINUE, Resultado: OBJ_NONE}
	ITER_ST     = &Status{Tipo: ITERACAO, Resultado: OBJ_NONE}
	SWITCH_ST   = &Status{Tipo: SWITCH, Resultado: OBJ_NONE}
	DECLARACAO_ST = &Status{Tipo: DECLARACAO,Resultado: OBJ_NONE}
	DEFINICAO_ST = &Status{Tipo: DEFINICAO,Resultado: OBJ_NONE}
)

type Status struct {
	Tipo      TipoStatus
	Resultado ObjetoBase
}
