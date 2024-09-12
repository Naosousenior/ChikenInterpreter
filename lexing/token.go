package lexing

import "fmt"

type TipoToken string

type Posicao struct {
	Linha  int
	Coluna int
}

type Token struct {
	Tipo  TipoToken
	Valor string
	Pos   Posicao
}

func (p Token) StringPos() string {
	return fmt.Sprintf("linha: %d, coluna: %d", p.Pos.Linha, p.Pos.Coluna)
}

func (p Token) IsTokenRecebe() bool {
	return p.Tipo == RECEBE || p.Tipo == SUB_RECEBE || p.Tipo == ADD_RECEBE || p.Tipo == MOD_RECEBE || p.Tipo == DIV_RECEBE || p.Tipo == MUL_RECEBE || p.Tipo == TIPO_RECEBE
}

func newToken(tipo TipoToken, valor byte) Token {
	return Token{Tipo: tipo, Valor: string(valor)}
}

func getIdentificador(valor string) TipoToken {
	var palavras_chave = map[string]TipoToken{
		"var":       VAR,
		"def":       DEFINITION,
		"run":       RUN,
		"class":     CLASS,
		"object":    OBJECT,
		"public":    PUBLICO,
		"protected": PROTEGIDO,
		"private":   PRIVADO,
		"super":     SUPER,
		"fun":       FUN,
		"repeat":    REPEAT,
		"iter":      ITER,
		"break":     BREAK,
		"continue":  CONTINUE,
		"if":        IF,
		"else":      ELSE,
		"switch":    SWITCH,
		"case":      CASE,
		"default":   DEFAULT,
		"true":      TRUE,
		"false":     FALSE,
		"none":      NONE,
		"return":    RETURN,
		"err":       ERROR,
		"try":       TRY,
		"except":    EXCEPT,
	}

	if token, teste := palavras_chave[valor]; teste {
		return token
	}
	return IDENTIFICADOR
}

const (
	RECEBE           = "="
	ADD_RECEBE       = "+="
	SUB_RECEBE       = "-="
	DIV_RECEBE       = "/="
	MUL_RECEBE       = "*="
	MOD_RECEBE       = "%="
	TIPO_RECEBE      = ":="
	PONTO_VIRGULA    = ";"
	DOIS_PONTO       = ":"
	TESTE_IS         = "::"
	ABRE_COLCHETE    = "["
	FECHA_COLCHETE   = "]"
	ABRE_CHAVE       = "{"
	FECHA_CHAVE      = "}"
	ABRE_PARENTESES  = "("
	FECHA_PARENTESES = ")"
	VIRGULA          = ","
	PONTO            = "."

	ADD            = "+"
	SUB            = "-"
	DIV            = "/"
	MUL            = "*"
	POTENCIA       = "**"
	RESTO          = "%"
	NEGACAO        = "!"
	COMPARACAO_IG  = "=="
	COMPARACAO_DIF = "!="
	MENOR_Q        = "<"
	MAIOR_Q        = ">"
	MENOR_IGUAL    = "<="
	MAIOR_IGUAL    = ">="
	AND            = "&"
	OR             = "|"
	XOR            = "||"
	SETA           = "->"

	STRING        = "STRING"
	COMENT        = "COMENT"
	NUM_REAL      = "REAL"
	NUM_INT       = "INT"
	IDENTIFICADOR = "IDENT"
	VAR           = "VAR"
	DEFINITION    = "DEFINITION"
	RUN           = "RUN"
	FUN           = "FUN"
	CLASS         = "CLASS"
	OBJECT        = "OBJECT"
	PUBLICO       = "PUBLICO"
	PROTEGIDO     = "PROTEGIDO"
	PRIVADO       = "PRIVADO"
	SUPER         = "SUPER"
	RETURN        = "RETURN"
	ERROR         = "ERROR"

	REPEAT   = "REPEAT"
	ITER     = "ITER"
	BREAK    = "BREAK"
	CONTINUE = "CONTINUE"
	IF       = "IF"
	ELSE     = "ELSE"
	SWITCH   = "SWITCH"
	CASE     = "CASE"
	DEFAULT  = "DEFAULT"
	TRY      = "TRY"
	EXCEPT   = "EXCEPT"

	TRUE  = "TRUE"
	FALSE = "FALSE"
	NONE  = "NONE"

	ILEGAL = "ILEGAL"
	FIM    = "FIM"
)
