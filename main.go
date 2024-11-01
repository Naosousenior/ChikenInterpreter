package main

import (
	"ChikenInterpreter/evaluation"
	"ChikenInterpreter/lexing"
	"ChikenInterpreter/objetos"
	"ChikenInterpreter/parsing"
	arv "ChikenInterpreter/parsing/arvore"
	"bufio"
	"fmt"
	"os"
)

var leitor_entrada = bufio.NewScanner(os.Stdin)

func ler_linha() string {
	if leitor_entrada.Scan() {
		return leitor_entrada.Text()
	}

	return "sem retorno"
}

func IteraTokens(texto string) {
	var lex *lexing.Lexico
	var tok lexing.Token
	lex = lexing.NewLexico(texto)
	tok = lex.GetToken()

	for tok.Tipo != lexing.FIM {
		fmt.Printf("Tipo: %s, Valor: %s, Posicao: %s\n", tok.Tipo, tok.Valor, tok.StringPos())
		tok = lex.GetToken()
	}
}

func IteraInstrucoes(texto string) {
	var lex *lexing.Lexico
	var analisador *parsing.Analisador

	lex = lexing.NewLexico(texto)
	analisador = parsing.NewAnalisador(lex)

	programa := analisador.AnalisaPrograma()

	erros := analisador.GetErros()

	if len(erros) > 0 {
		for _, msg := range erros {
			fmt.Println(msg)
		}

	} else {
		fmt.Println("Nenhum erro encontrado")
	}

	fmt.Println(programa.GetInformacao())

}

func ModoInterativo() {
	var comando string
	var lexico *lexing.Lexico
	var analisador *parsing.Analisador
	var erros []string
	var programa *arv.Programa
	ambiente := objetos.NewAmbiente()
	ambiente.CriaVar("bff",evaluation.NewBff())

	fmt.Println("Its are an command line interpreter!!")

	for {

		fmt.Print(">> ")
		comando = ler_linha()

		if comando == "exit" {
			break
		}

		if len(comando) > 4 {
			if comando[0:3] == "run" {
				arquivo := comando[4:]
				fmt.Println(arquivo)
				texto := LerArquivo(arquivo + ".txt")

				if texto != "" {
					lexico = lexing.NewLexico(texto)
					analisador = parsing.NewAnalisador(lexico)
					programa = analisador.AnalisaPrograma()

					erros := analisador.GetErros()

					if len(erros) > 0 {
						for _, err := range erros {
							fmt.Println(err)
						}

						continue
					}

					resultado := evaluation.Avaliar(programa, ambiente)
					if resultado != nil {
						fmt.Println(resultado.Inspecionar())
					} else {
						fmt.Println("Entrada invalida, tente outra coisa")
					}
				}

				continue
			}
		}

		if comando[len(comando)-1] != ';' {
			comando += ";"
		}

		lexico = lexing.NewLexico(comando)
		analisador = parsing.NewAnalisador(lexico)
		programa = analisador.AnalisaPrograma()

		erros = analisador.GetErros()
		if len(erros) > 0 {
			for _, err := range erros {
				fmt.Println(err)
			}
			continue
		}

		resultado := evaluation.Avaliar(programa, ambiente)
		if resultado == nil {
			fmt.Println("Algo deu errado")
			continue
		}
		fmt.Println(resultado.Inspecionar())
	}
}

func LerArquivo(file string) string {
	texto_b, erro := os.ReadFile(file)
	texto := string(texto_b)

	if erro != nil {
		return ""
	}

	return texto
}

func main() {

	ModoInterativo()

	texto := LerArquivo("teste.txt")

	IteraInstrucoes(texto)
}
