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

type Conjunto map[string]struct{}

func (c Conjunto) Add(elemento string) {
	c[elemento] = struct{}{}
}

func (c Conjunto) Remove(elemento string) {
	delete(c, elemento)
}

func (c Conjunto) Tem(elemento string) bool {
	_, ok := c[elemento]

	return ok
}

func (c Conjunto) Copiar(conjunto Conjunto) {
	for chave := range conjunto {
		c[chave] = struct{}{}
	}
}
