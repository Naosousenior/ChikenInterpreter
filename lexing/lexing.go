package lexing

import (
	"ChikenInterpreter/ferramentas"
)

type Lexico struct {
	texto   string
	posicao int
	proximo int
	letra   byte

	linha  int
	coluna int
}

func (l *Lexico) leLetra() {
	if l.proximo >= len(l.texto) {
		l.letra = 0
	} else {
		l.letra = l.texto[l.proximo]
	}
	l.posicao = l.proximo
	l.proximo += 1

	l.coluna += 1
}

func (l *Lexico) espiaLetra(qtd int) byte {
	if l.proximo+qtd >= len(l.texto) {
		return 0
	}
	return l.texto[l.proximo+qtd]
}

func (l *Lexico) pula_espaco_branco() {
	for l.letra == ' ' || l.letra == '\n' || l.letra == '\t' || l.letra == '\r' {
		if l.letra == '\n' {
			l.coluna = 0
			l.linha++
		}
		l.leLetra()
	}
}

func (l *Lexico) le_palavra() string {
	resultado := ""
	for {
		resultado += string(l.letra)
		if !(ferramentas.ELetra(l.espiaLetra(0)) || l.espiaLetra(0) == '_' || ferramentas.ENumero( l.espiaLetra(0) ) ) {
			break
		}
		l.leLetra()
	}

	return resultado
}

func (l *Lexico) le_numero() (string, TipoToken) {
	resultado := ""
	tipoToken := NUM_INT
	qtd_ponto := 0
	for {
		resultado += string(l.letra)

		teste := ferramentas.ENumero(l.espiaLetra(0)) || l.espiaLetra(0) == '.'
		if !teste {
			break
		} else if l.espiaLetra(0) == '.' {
			if qtd_ponto > 0 {
				break
			}
			if !(ferramentas.ENumero(l.espiaLetra(1))) {
				break
			} else {
				tipoToken = NUM_REAL
				qtd_ponto++
			}
		}

		l.leLetra()
	}

	return resultado, TipoToken(tipoToken)
}

func (l *Lexico) le_string() string {
	resultado := ""
	l.leLetra()
	for !(l.letra == '\'' || l.letra == 0) {
		resultado += string(l.letra)
		l.leLetra()
	}

	return resultado
}

func (l *Lexico) le_coment() string {
	resultado := ""
	l.leLetra()
	for l.letra != '"' {
		resultado += string(l.letra)
		l.leLetra()
	}

	return resultado
}

func NewLexico(texto string) *Lexico {
	l := &Lexico{texto: texto, proximo: 0, linha: 1, coluna: 0}
	l.leLetra()

	return l
}

func (l *Lexico) GetToken() Token {
	var tok Token

	l.pula_espaco_branco()

	posicao := Posicao{Linha: l.linha, Coluna: l.coluna}

	switch l.letra {
	case ';':
		tok = newToken(PONTO_VIRGULA, l.letra)
	case ':':
		if l.espiaLetra(0) == '=' {
			tok.Tipo = TIPO_RECEBE
			tok.Valor = ":="
			l.leLetra()
		} else if l.espiaLetra(0) == ':' {
			tok.Tipo = TESTE_IS
			tok.Valor = "::"
			l.leLetra()
		} else {
			tok = newToken(DOIS_PONTO, l.letra)
		}
	case '.':
		tok = newToken(PONTO, l.letra)
	case ',':
		tok = newToken(VIRGULA, l.letra)
	case '=':
		if l.espiaLetra(0) == '=' {
			tok.Tipo = COMPARACAO_IG
			tok.Valor = "=="
			l.leLetra()
		} else {
			tok = newToken(RECEBE, l.letra)
		}
	case '!':
		if l.espiaLetra(0) == '=' {
			tok.Tipo = COMPARACAO_DIF
			tok.Valor = "!="
			l.leLetra()
		} else {
			tok = newToken(NEGACAO, l.letra)
		}

	case '<':
		if l.espiaLetra(0) == '=' {
			tok.Tipo = MENOR_IGUAL
			tok.Valor = "<="
			l.leLetra()
		} else {
			tok = newToken(MENOR_Q, l.letra)
		}
	case '>':
		if l.espiaLetra(0) == '=' {
			tok.Tipo = MAIOR_IGUAL
			tok.Valor = ">="
			l.leLetra()
		} else {
			tok = newToken(MAIOR_Q, l.letra)
		}
	case '*':
		if l.espiaLetra(0) == '*' {
			tok.Tipo = POTENCIA
			tok.Valor = "**"
			l.leLetra()
		} else if l.espiaLetra(0) == '=' {
			tok.Tipo = MUL_RECEBE
			tok.Valor = "*="
			l.leLetra()
		} else {
			tok = newToken(MUL, l.letra)
		}
	case '[':
		tok = newToken(ABRE_COLCHETE, l.letra)
	case ']':
		tok = newToken(FECHA_COLCHETE, l.letra)
	case '(':
		tok = newToken(ABRE_PARENTESES, l.letra)
	case ')':
		tok = newToken(FECHA_PARENTESES, l.letra)
	case '{':
		tok = newToken(ABRE_CHAVE, l.letra)
	case '}':
		tok = newToken(FECHA_CHAVE, l.letra)
	case '+':
		if l.espiaLetra(0) == '=' {
			tok.Tipo = ADD_RECEBE
			tok.Valor = "+="
			l.leLetra()
		} else {
			tok = newToken(ADD, l.letra)
		}
	case '-':
		if l.espiaLetra(0) == '=' {
			tok.Tipo = SUB_RECEBE
			tok.Valor = "-="
			l.leLetra()
		} else if l.espiaLetra(0) == '>' {
			tok.Tipo = SETA
			tok.Valor = "->"
			l.leLetra()
		} else {
			tok = newToken(SUB, l.letra)
		}
	case '/':
		if l.espiaLetra(0) == '=' {
			tok.Tipo = DIV_RECEBE
			tok.Valor = "/="
			l.leLetra()
		} else {
			tok = newToken(DIV, l.letra)
		}
	case '%':
		if l.espiaLetra(0) == '=' {
			tok.Tipo = MOD_RECEBE
			tok.Valor = "%="
			l.leLetra()
		} else {
			tok = newToken(RESTO, l.letra)
		}
	case '&':
		tok = newToken(AND, l.letra)
	case '|':
		if l.espiaLetra(0) == '|' {
			tok.Tipo = XOR
			tok.Valor = "||"
			l.leLetra()
		} else {
			tok = newToken(OR, l.letra)
		}
	case '\'':
		tok.Tipo = STRING
		tok.Valor = l.le_string()
	case '"':
		tok.Tipo = COMENT
		tok.Valor = l.le_coment()
	case 0:
		tok.Tipo = FIM
		tok.Valor = ""

	default:
		if ferramentas.ELetra(l.letra) || l.letra == '_' {
			palavra := l.le_palavra()
			tok.Tipo = getIdentificador(palavra)
			tok.Valor = palavra
		} else if ferramentas.ENumero(l.letra) {
			numero, tipo := l.le_numero()
			tok.Tipo = tipo
			tok.Valor = numero
		} else {
			tok = newToken(ILEGAL, l.letra)
		}
	}

	tok.Pos = posicao

	if tok.Tipo != FIM {
		l.leLetra()
	}
	return tok
}
