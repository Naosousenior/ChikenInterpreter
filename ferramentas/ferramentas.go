package ferramentas

func ELetra(letra byte) bool {
	if (letra >= 'a' && letra <= 'z') || (letra >= 'A' && letra <= 'Z') {
		return true
	}

	return false
}

func ENumero(letra byte) bool {
	if letra >= '0' && letra <= '9' {
		return true
	}

	return false
}

func GetIdentacao(profundidade int) string {
	esp := "    "
	identacao := ""
	i := 0

	for i < profundidade {
		identacao += esp
		i++
	}

	return identacao
}