package objetos

type Ambiente struct {
	variaveis map[string]ObjetoBase
	definicoes map[string]ObjetoBase
	externo   *Ambiente
	Classe    *Classe
	Objeto    *ObjetoUser
}

func NewAmbiente() *Ambiente {


	return &Ambiente{externo: nil,variaveis: make(map[string]ObjetoBase),definicoes: make(map[string]ObjetoBase)}
}

func NewAmbienteInterno(amb *Ambiente) *Ambiente {
	novo := NewAmbiente()

	novo.externo = amb
	return novo
}

func (a *Ambiente) existeExterno() bool {
	return a.externo != nil
}

func (a *Ambiente) existeAqui(nome string) bool {
	_,ok1 := a.variaveis[nome]
	_,ok2 := a.definicoes[nome]

	return ok1 || ok2
}

//fragmentei essa tarefa em 3 funcoes
func (a *Ambiente) existe(nome string) bool {
	if a.existeAqui(nome) {
		return true
	}

	if !a.existeExterno() {
		return false
	}

	return a.externo.existeAqui(nome) 
}

func (a *Ambiente) eAtribuivel(nome string) bool {
	_,ok := a.definicoes[nome]

	if ok { //se existe nas definicoes do escopo atual, ja nao pode ser reatribuido
		return false
	}

	if a.existeExterno() { //se nao existe, vai verificar no escopo exterior
		return a.externo.eAtribuivel(nome)
	}

	return true //se nao existe um escopo exterior, pode ser reatribuido
	//nao verifica e existencia da variavel
}

func (a *Ambiente) get_var(nome string) ObjetoBase {
	//primeiro, tenta pegar entre as variaveis
	if res_var,ok_var := a.variaveis[nome]; ok_var {
		return res_var
	}

	//depois, entre as definicoes
	if res_def,ok_def := a.definicoes[nome]; ok_def {
		return res_def
	}

	//como temos certeza que existe, vamos tentar no ambiente exterior
	return a.externo.get_var(nome)
}

func (a *Ambiente) GetVar(nome string) (ObjetoBase, bool) {
	if !a.existe(nome) { //verifica a existencia da variavel
		return nil,false
	}

	//se passou do teste anterior, e porque existe, nao tem porque temer um erro
	return a.get_var(nome),true
}

func (a *Ambiente) AddArgs(nome string, valor ObjetoBase) {
	a.variaveis[nome] = valor
}

func (a *Ambiente) CriaVar(nome string, variavel ObjetoBase) ObjetoBase {
	if !a.existe(nome) {
		a.variaveis[nome] = variavel
		return OBJ_NONE
	}

	return &ObjExcessao{Mensagem: "Variavel " + nome + " ja existente."}
}

func (a *Ambiente) DefVar(nome string, variavel ObjetoBase) bool {
	if a.existe(nome) {
		return false
	}

	a.definicoes[nome] = variavel
	return true
}

func (a *Ambiente) set_var(ref string,variavel ObjetoBase) {
	if a.existeAqui(ref) {
		a.variaveis[ref] = variavel
		return
	}

	a.externo.set_var(ref,variavel) //se essa funcao foi chamada, e porque a referencia existe
	//se ela existe, mas não existe no escopo atual, deve estar no escopo exterior, que por sua vez, deve existir
}

func (a *Ambiente) SetVar(ref string, variavel ObjetoBase) bool {
	if a.existe(ref) && a.eAtribuivel(ref) { //se passou, é porque existe e pode ser modificado
		a.set_var(ref,variavel)
		return true
	}

	return false
}
