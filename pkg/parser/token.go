package parser

// TokenType es un alias para string, usado para representar el tipo de token.
type TokenType string

// Token representa un token léxico, con su tipo y el literal original.
type Token struct {
	Type    TokenType
	Literal string
}

const (
	// Tipos de tokens especiales
	ILLEGAL TokenType = "ILLEGAL" // Carácter o secuencia no reconocida
	EOF     TokenType = "EOF"     // Fin de archivo / entrada

	// Identificadores y literales
	IDENT  TokenType = "IDENT"  // Nombres de comandos, variables, etc. (ej. ls, my_var)
	INT    TokenType = "INT"    // Números enteros (ej. 123)
	STRING TokenType = "STRING" // Cadenas de texto (ej. "hello world", 'otra cadena')

	// Operadores
	ASSIGN   TokenType = "="   // Asignación (let x = 10)
	PIPE     TokenType = "|"   // Pipe de comandos
	EQ       TokenType = "=="  // Igualdad
	NEQ      TokenType = "!="  // Desigualdad
	GT       TokenType = ">"   // Mayor que
	LT       TokenType = "<"   // Menor que
	GTE      TokenType = ">="  // Mayor o igual que
	LTE      TokenType = "<="  // Menor o igual que
	DOT      TokenType = "."   // Acceso a campo (ej. .data.name)

	// Delimitadores
	COMMA    TokenType = ","   // Coma (para separar argumentos, etc.)
	SEMICOLON TokenType = ";"  // Punto y coma (para separar statements, si se implementa)
	LPAREN   TokenType = "("   // Paréntesis izquierdo
	RPAREN   TokenType = ")"   // Paréntesis derecho
	LBRACE   TokenType = "{"   // Llave izquierda (para bloques de código)
	RBRACE   TokenType = "}"   // Llave derecha
	LBRACKET TokenType = "["   // Corchete izquierdo (para arrays o indexing)
	RBRACKET TokenType = "]"   // Corchete derecho

	// Palabras clave
	LET    TokenType = "LET"    // 'let' keyword
	CD     TokenType = "CD"     // 'cd' keyword
	VARS   TokenType = "VARS"   // 'vars' keyword
	EXIT   TokenType = "EXIT"   // 'exit' keyword
	GET    TokenType = "GET"    // 'get' built-in
	WHERE  TokenType = "WHERE"  // 'where' built-in
	SELECT TokenType = "SELECT" // 'select' built-in
	IF     TokenType = "IF"     // 'if' keyword
	ELSE   TokenType = "ELSE"   // 'else' keyword
	FOR    TokenType = "FOR"    // 'for' keyword
	DEF    TokenType = "DEF"    // 'def' keyword (para definir funciones)
	TRUE   TokenType = "TRUE"   // 'true' boolean literal
	FALSE  TokenType = "FALSE"  // 'false' boolean literal
)

// keywords es un mapa para buscar si una cadena es una palabra clave.
var keywords = map[string]TokenType{
	"let":    LET,
	"cd":     CD,
	"vars":   VARS,
	"exit":   EXIT,
	"get":    GET,
	"where":  WHERE,
	"select": SELECT,
	"if":     IF,
	"else":   ELSE,
	"for":    FOR,
	"def":    DEF,
	"true":   TRUE,
	"false":  FALSE,
}

// LookupIdent verifica si el identificador dado es una palabra clave
// o un identificador de usuario.
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}
