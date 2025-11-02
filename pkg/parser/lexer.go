package parser

import "unicode"

// Lexer se encarga de tokenizar la entrada.
type Lexer struct {
	input        string
	position     int  // posición actual en la entrada (apunta al carácter actual)
	readPosition int  // próxima posición a leer (después del carácter actual)
	ch           rune // carácter actual bajo inspección
}

// NewLexer crea una nueva instancia de Lexer.
func NewLexer(input string) *Lexer {
	l := &Lexer{input: input}
	l.readChar() // Inicializa la posición y el primer carácter
	return l
}

// readChar avanza la posición en el input y lee el siguiente carácter.
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // ASCII para "NUL", indica EOF
	} else {
		l.ch = rune(l.input[l.readPosition])
	}
	l.position = l.readPosition
	l.readPosition++
}

// peekChar devuelve el siguiente carácter sin avanzar la posición.
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	return rune(l.input[l.readPosition])
}

// NextToken tokeniza el siguiente fragmento de la entrada.
func (l *Lexer) NextToken() Token {
	var tok Token

	l.skipWhitespace()

	switch l.ch {
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: EQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ASSIGN, l.ch)
		}
	case '|':
		tok = newToken(PIPE, l.ch)
	case ';':
		tok = newToken(SEMICOLON, l.ch)
	case '(':
		tok = newToken(LPAREN, l.ch)
	case ')':
		tok = newToken(RPAREN, l.ch)
	case '{':
		tok = newToken(LBRACE, l.ch)
	case '}':
		tok = newToken(RBRACE, l.ch)
	case '[':
		tok = newToken(LBRACKET, l.ch)
	case ']':
		tok = newToken(RBRACKET, l.ch)
	case ',':
		tok = newToken(COMMA, l.ch)
	case '.':
		// Un punto puede ser un token por sí mismo (para `get`) o parte de un identificador.
		// Si está seguido por espacio o nada, es un token. Si no, será parte de un identificador.
		if !isIdentifierChar(l.peekChar()) {
			tok = newToken(DOT, l.ch)
		} else {
			// Cae en el caso default para ser leído como parte de un identificador.
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal)
			return tok
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: GTE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(GT, l.ch)
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: LTE, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(LT, l.ch)
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = Token{Type: NEQ, Literal: string(ch) + string(l.ch)}
		} else {
			tok = newToken(ILLEGAL, l.ch) // '!' solo no es válido por ahora
		}
	case '"', '\'':
		tok.Type = STRING
		tok.Literal = l.readString(l.ch)
	case 0:
		tok.Literal = ""
		tok.Type = EOF
	default:
		if isIdentifierChar(l.ch) {
			tok.Literal = l.readIdentifier()
			tok.Type = LookupIdent(tok.Literal) // Verifica si es una palabra clave
			return tok
		} else if isDigit(l.ch) {
			tok.Type = INT
			tok.Literal = l.readNumber()
			return tok
		} else {
			tok = newToken(ILLEGAL, l.ch)
		}
	}

	l.readChar() // Avanza al siguiente carácter después de tokenizar
	return tok
}

// newToken es una función auxiliar para crear tokens.
func newToken(tokenType TokenType, ch rune) Token {
	return Token{Type: tokenType, Literal: string(ch)}
}

// skipWhitespace avanza el lexer pasando los espacios en blanco.
func (l *Lexer) skipWhitespace() {
	for l.ch == ' ' || l.ch == '\t' || l.ch == '\n' || l.ch == '\r' {
		l.readChar()
	}
}

// readIdentifier lee un identificador (o palabra clave) hasta que encuentra un no-letra/dígito.
func (l *Lexer) readIdentifier() string {
	position := l.position
	for isIdentifierChar(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readNumber lee un número (entero por ahora) hasta que encuentra un no-dígito.
func (l *Lexer) readNumber() string {
	position := l.position
	for isDigit(l.ch) {
		l.readChar()
	}
	return l.input[position:l.position]
}

// readString lee una cadena entre comillas (simples o dobles).
func (l *Lexer) readString(quote rune) string {
	position := l.position + 1 // Ignora la comilla de apertura
	for {
		l.readChar()
		if l.ch == quote || l.ch == 0 { // Si encontramos la comilla de cierre o EOF
			break
		}
	}
	return l.input[position:l.position]
}

// isIdentifierChar verifica si el rune es un carácter válido para un identificador o argumento.
// CORREGIDO: Esta es la nueva definición, mucho más permisiva.
func isIdentifierChar(ch rune) bool {
	switch ch {
	case ' ', '\t', '\n', '\r', '|', '=', '(', ')', '{', '}', '[', ']', '"', '\'', 0:
		return false
	default:
		return true
	}
}

// isDigit verifica si el rune es un dígito.
func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}
