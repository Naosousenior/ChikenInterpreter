
# APRENDENDO A DESENVOLVER UM INTERPRETADOR  

Este interpretador foi desenvolvido com auxílio do livro "Writing a interpreter in Go", de Thorsten Ball.
O objeto no início era realmente só criar um interpretador, mas acabei me apegando a linguagem que eu mesmo comecei a desenvolver. A versão atual é meio instável, no entanto, alguns recursos já são utilizáveis.  
Agora, estou tentando desenvolver uma _linguagem orientada a objetos_, e como tive péssimas experiências com POO em outras linguagens de script, como JS, resolvi que seria uma orientação a objetos robusta. No momento, ainda estou tentando implementar essa parte, mas os demais recursos da linguagem já estão prontos.

# RECURSOS DA LINGUAGEM
Agora, vamos dar uma olhada nos recursos que eu já desenvolvi para a linguagem Chiken, a maioria deles com auxílio do livro referido acima (aliás, recomendo fortemente a quem está interessado neste conteúdo).

## Expressoes
Em Chiken, expressões são muito semelhantes as outras linguagens. Algo como `a+b`, ou `a*b`. A seguir a tabela completa de operações já implementadas em Chiken:

Simbolo | Conceito | Exemplo
------- | -------- | -------
  +| Adição  | `a+b`
  -|Subtração| `a-b`
  *|Multiplicação|`a*b`
  \\|Divisão|`a\b`
  % |Resto da divisão|`a%b`
  **|Potenciação|`a**b`
  :|Verificador de pertencimento (verifica se o elemento `b` pertence a `a`)|`a:b`
  ::|Verificador de tipo (verifica se `b` é do tipo `a`)|`a::b`
  !|Negação|`!a`
  ==|Igualdade|`a==b`
  !=|Diferença|`a!=b`
  <|Menor que|`a<b`
  \>|Maior que|`a>b`
  <=, >=|Menor ou igual e maior ou igual|`a<=b,a>=b`
  &|And (e)|`a & b`
  \||Or (ou)|`a \| b`
  \|\| |Xor (ou exclusivo)|`a \|\|b`
  
Por fim, devo lembrar que nem internamente, nem externamente, `.` nunuca será considerado um operador. Afinal, você não pode colocar uma expressão a direita de `.`, você **sempre** terá de colocar um identificador, pois para a linguagem, `.` é um acesso a um atributo.
## Declaração de variáveis
Declarar variáveis em Chiken é tão simples como qualquer outra linguagem: basta usar a palavra chave `var`:
```
var <nome> = <expressao objeto>;
```

Aliás, sim, em Chiken o ponto e vírgula é obrigatório no final dos comandos (exceto, obviamente, no final de blocos de comando). Mas a declaração de variáveis não para por aqui. Podemos também declarar uma variável sem seu valor inicial.
```
var <nome>;
```
Implicitamente, o valor atribuido a essa variável sera `none`. Em Chiken, `none` é um objeto vazio, ou o ponteiro nulo. Entenda como você quiser. Mas `none` possui certos métodos, e pode ser tratado em funções, por isso prefiro a denominação "objeto vazio". A expressão acima é equivalente a:
```
var <nome> = none;
```
O ponto mais interessanta talvez seja o fato de que você pode declarar uma lista de varíaveis, com seus valores, separando com vírgula, em uma única instrução `var`. Mais ou menos assim:
```
var <nome1> = <expressao1>,<nome2> = <expressao2>,<nome3>,<nome4> = <expressao3>;
```

Obviamente, é mais elegante fazer quebras de linhas.
```
var <nome1> = <expressao1>,
<nome2> = <expressao2>,
<nome3>,
<nome4> = <expressao3>;
```

### Funções
Funções são um tipo especial de variável. Afinal, elas podem carregar código dentro de si. Em ChikenLang, todas as função são de ordem superior, ou seja, ela podem ser passadas como parâmetro de outras funções, ou podem ser retornadas por funções.
Para definir uma função, basta utilizar a palavra chave `fn`:
```
var funcao = fn(<parametros>) {
	<comandos>
};
```
Para retornar um valor em uma função, use a instrução `return`:
```
var funcao = fn(<parametros>) {
	return <valor>;
};
```
## Estruturas de blocos
Talvez essa seja a parte mais interessante. Existem 4 estruturas na linguagem até agora (ainda quero implementar algumas outras). São elas:
### If
O bloco `if` em Chiken não é nada diferente de outras linguagens:
```
if(<condição>){
    <comandos...>
}
```
e obviamente, também é possível adiconar o else:
```
if(<condição>){
    <comandos...>
} else {
    <outros comandos...>
}
```
além de `if-else` aninhados:
```
if(<condição1>){
    <comandos...>
} else if(<condição2) {
    <comandos...>
} else if
...
} else {
    <comandos...>
}
```
Uma característica importante do comando `if`, é que ele também pode ser usado como uma expressão, ou seja, ele pode gerar um valor. Isso significa que você também pode usar ele para atribuir um valor a uma variável.
```
var <nome> = if(<condicao){<comandos>}
```

Se a condição for falsa e não houver um else, o valor retornado será `none`.
### Repeat
Esse é o comando básico de repetição. Ele substitui a necessidade de `while` e `do-while`, possibilitando o processo de looping. Sua estrutura pode ser descrita assim:
```
repeat(<condicao1) {
    <comandos>
}(<condicao2)
```

Simples assim. Em resumo, `repeat` vai testar a condição 1. Se ela for verdadeira, ele executa o bloco de comandos. Depois, ele verifica a condição 2. Se ela for verdadeira, ele reinicia o looping, testando novamente a condição 1. Ambas as condições, 1 e 2, são opcionais, ou seja, você pode usar o repeat assim
```
repeat(<condicao1) {
    <comandos>
}(<condicao2)
```
Assim:
```
repeat(<condicao1) {
    <comandos>
}
```
Assim:
```
repeat {
    <comandos>
}(<condicao2)
```
E até assim:
```
repeat {
    <comandos>
}
```
As condições omitidas terão valor `true` por padrão. As instruções `break`, `continue` e `return` modificam o ciclo natural do `repeat`. `break` e `return` saem do looping, enquanto `continue` pula do comando do bloco para o próximo teste da condição 2.  
Finalmente, `repeat` assim como `if` também pode retornar um objeto, e ser usado como expressões. No caso, `repeat` retorna uma lista, com os valores produzidos pelo bloco de código a cada iteração. Então, por exemplo, o código a seguir produz um array de 0 a 9:
```
var inicio = -1;
var array = repeat(inicio<9) {
    inicio +=1 ;
    inicio;
};
```
### Iter
Instrução de iteração em objetos iteráveis. Por exemplo, um array é iterável. Um exemplo está a seguir:
```
iter(<iterador> : <iteravel>) {
    <comandos>
}
```
Arrays e dicts internos da linguagem são exemplos de objetos iteráveis.
### Switch
Instrução de seleção de casos. Uma expressão pode ser avaliada para gerar um valor que será utilizado para avançar para cada caso:
```
switch(<expressao>) {
    <caso1> {
	    <codigo>
    }
    <caso2> {
	    <codigo>
	}
	...
}
```
Em teoria, era pra ser bem leve. Mas provavelmente será lento pra caramba, recomendo fortemente você buscar algum método alternativo.

## Orientação a objetos
Chiken Lang é uma linguagem orientada a objetos, e isso foi um desafio pessoal. Como tal, ela implementa encapsulamento, herança múltipla, e obviamente, polimorfismo e instanciação de objetos. Antes de mais nada, em Chiken Lang um objeto pode ser criado sem necessariamente ter uma classe, por meio da palavra chave `object`:
```
var objeto = object {
	attrib1 = <value1>,
	attrib2 = <value2>
};
```
se você fizer assim, o objeto será basicamente uma estrutura. Funções definidas dentro dele serão métodos, no entanto, será impossível fazer encapsulamento, além de definir um construtor. Na verdade, um objeto assim será único, a menos que você o copie, e não vejo porque fazer isso. Este recurso é útil para o caso de você querer um objeto único que faça alguma coisa, sem usar Singleton.  
Para criar uma classe, use a palavra chave `class`:
```
var classe = class(<lista de super classes>) {
	<atributos>
};
```
Os atributos serão definidos da mesma forma que no caso do comando `object`, no entanto, modificadores de encapsulamento podem ser usados:
```
var classe = class(<lista de super classes>) {
	public
	<atributos publicos>
	
	protected
	<atributos protegidos>
	
	private
	<atributos privados>
	
	class
	<atributos de classe>
};
```
atributos de classe não fazem parte do encapsulamento, obviamente, eles são exatamente o que foi descrito: um atributo que pertence a classe, e que portanto, não será atribuído aos objetos instanciados por aquela classe.  
Atributos públicos podem ser acessados por qualquer um.  
Atributos protegidos podem ser acessados pela própria classe e por classes filhas.  
Atributos privados só podem ser acessados pela própria classe.  
O modificador padrão é `public`.  
Um método é definido ao se criar um atributo que guarda uma função:
```
var classe = class(<lista de super classes>) {
	metodo = fn(<parametros>) {
		<comandos>
	},
};
```
Um método especial é o construtor, ele deve possuir o nome `new_object` e ser público:
```
var classe = class(<lista de super classes>) {
	new_object = fn(<parametros>) {
		<comandos>
	},
};
```
Como foi citado, em Chiken Lang existe herança múltipla, e sobre isso, é importante ter em mente que a ordem de prioridade dos atributos é da primeira classe passada para a última. Ou seja, supondo que duas classes herdadas tenham um atributo de mesmo nome, a classe mais a esquerda será priorizada em relação a mais a direita.