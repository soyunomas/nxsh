package parser

import (
	"fmt"
	"strings"
)

// Parser toma un lexer y construye un AST.
type Parser struct {
	l      *Lexer
	errors []string

	curToken  Token
	peekToken Token
}

// NewParser crea una nueva instancia del Parser.
func NewParser(l *Lexer) *Parser {
	p := &Parser{
		l:      l,
		errors: []string{},
	}
	p.nextToken()
	p.nextToken()
	return p
}

func (p *Parser) Errors() []string {
	return p.errors
}

func (p *Parser) peekError(t TokenType) {
	msg := fmt.Sprintf("se esperaba que el siguiente token fuera %s, pero se obtuvo %s",
		t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

func (p *Parser) nextToken() {
	p.curToken = p.peekToken
	p.peekToken = p.l.NextToken()
}

func (p *Parser) ParseProgram() *Program {
	program := &Program{}
	program.Statements = []Statement{}

	for p.curToken.Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	return program
}

func (p *Parser) parseStatement() Statement {
	switch p.curToken.Type {
	case LET:
		return p.parseLetStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *LetStatement {
	stmt := &LetStatement{Token: p.curToken}

	if !p.expectPeek(IDENT) {
		return nil
	}

	stmt.Name = &Identifier{Token: p.curToken, Value: p.curToken.Literal}

	if !p.expectPeek(ASSIGN) {
		return nil
	}

	p.nextToken()

	stmt.Value = p.parseExpression()

	// Es crucial verificar si la expresión fue parseada correctamente.
	if stmt.Value == nil {
		msg := fmt.Sprintf("no se encontró una expresión válida después de '=' en la declaración let")
		p.errors = append(p.errors, msg)
		return nil
	}

	return stmt
}

func (p *Parser) parseExpressionStatement() *ExpressionStatement {
	stmt := &ExpressionStatement{Token: p.curToken}
	stmt.Expression = p.parseExpression()

	if stmt.Expression == nil {
		return nil
	}

	return stmt
}

// ¡FUNCIÓN CORREGIDA!
// Esta es la lógica principal que soluciona el bug.
func (p *Parser) parseExpression() Expression {
	var left Expression

	// Decidimos qué tipo de expresión estamos viendo.
	if p.isCommandStartToken() {
		// Parece un comando, como `ls -l` o `curl ...`
		left = p.parseCommandExpression()
	} else {
		// Es otra cosa, como un StringLiteral `"hola"` o un futuro número.
		left = p.parsePrimaryExpression()
	}

	// Después de parsear la parte izquierda, comprobamos si le sigue un pipe.
	if p.peekTokenIs(PIPE) {
		if left == nil {
			p.errors = append(p.errors, "expresión inválida antes del pipe '|'")
			return nil
		}
		p.nextToken()
		pipeToken := p.curToken
		p.nextToken() // Avanzamos al inicio de la expresión derecha
		right := p.parseExpression()

		if right == nil {
			p.errors = append(p.errors, "expresión vacía o inválida después del pipe '|'")
			return nil
		}

		return &PipelineExpression{Token: pipeToken, Left: left, Right: right}
	}

	return left
}

func (p *Parser) parseCommandExpression() Expression {
	if !p.isCommandStartToken() {
		return nil
	}

	cmd := &CommandExpression{
		Token: p.curToken,
		Name:  &Identifier{Token: p.curToken, Value: p.curToken.Literal},
		Args:  []Expression{},
	}

	for !p.peekTokenIs(PIPE) && !p.peekTokenIs(EOF) {
		p.nextToken()
		arg := p.parsePrimaryExpression()
		if arg != nil {
			cmd.Args = append(cmd.Args, arg)
		}
	}

	return cmd
}

func (p *Parser) isCommandStartToken() bool {
	return p.curToken.Type == IDENT || p.curToken.Type == GET ||
		p.curToken.Type == WHERE || p.curToken.Type == SELECT ||
		p.curToken.Type == CD || p.curToken.Type == VARS || p.curToken.Type == EXIT
}

// parsePrimaryExpression parsea los componentes básicos de un comando.
func (p *Parser) parsePrimaryExpression() Expression {
	switch p.curToken.Type {
	case IDENT, GET, WHERE, SELECT, INT, TRUE, FALSE, EQ, NEQ, GT, LT, GTE, LTE:
		return &Identifier{Token: p.curToken, Value: p.curToken.Literal}
	case STRING:
		return &StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
	case DOT:
		return p.parsePathExpression()
	default:
		// Añadimos un error si no es una expresión que conocemos.
		p.errors = append(p.errors, fmt.Sprintf("no se pudo parsear la expresión que empieza con '%s'", p.curToken.Literal))
		return nil
	}
}

func (p *Parser) parsePathExpression() Expression {
	startToken := p.curToken
	var path strings.Builder
	path.WriteString(p.curToken.Literal)

	for {
		if p.curTokenIs(DOT) && p.peekTokenIs(IDENT) {
			p.nextToken()
			path.WriteString(p.curToken.Literal)
		} else if p.curTokenIs(IDENT) && p.peekTokenIs(DOT) {
			p.nextToken()
			path.WriteString(p.curToken.Literal)
		} else {
			break
		}
	}

	return &Identifier{Token: startToken, Value: path.String()}
}

func (p *Parser) curTokenIs(t TokenType) bool {
	return p.curToken.Type == t
}

func (p *Parser) peekTokenIs(t TokenType) bool {
	return p.peekToken.Type == t
}

func (p *Parser) expectPeek(t TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}
