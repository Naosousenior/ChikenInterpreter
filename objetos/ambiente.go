package objetos

type Ambiente struct {
	variaveis map[string]ObjetoBase
	externo   *Ambiente
	Classe    *Classe
	Objeto    *ObjetoUser
}

func NewAmbiente() *Ambiente {
	vars := make(map[string]ObjetoBase)

	return &Ambiente{variaveis: vars}
}

func NewAmbienteInterno(amb *Ambiente) *Ambiente {
	novo := NewAmbiente()

	novo.externo = amb
	return novo
}

func (a *Ambiente) GetVar(nome string) (ObjetoBase, bool) {
	res, ok := a.variaveis[nome]

	if !ok && a.externo != nil {
		res, ok = a.externo.GetVar(nome)
	}

	return res, ok
}

func (a *Ambiente) AddArgs(nome string, valor ObjetoBase) {
	a.variaveis[nome] = valor
}

func (a *Ambiente) CriaVar(nome string, variavel ObjetoBase) ObjetoBase {
	if _, ok := a.GetVar(nome); !ok {
		a.variaveis[nome] = variavel
		return OBJ_NONE
	}

	return &ObjExcessao{Mensagem: "Variavel " + nome + " ja existente."}
}

func (a *Ambiente) SetVar(ref string, variavel ObjetoBase) bool {
	if a.externo != nil {
		if a.externo.SetVar(ref, variavel) {
			return true
		} else if _, ok := a.GetVar(ref); ok {
			a.variaveis[ref] = variavel
			return true
		}

		return false
	} else if _, ok := a.GetVar(ref); ok {
		a.variaveis[ref] = variavel
		return true
	}

	return false
}
