// pkg/parser/ast.go
package parser

import (
	"bytes"
	// La línea "github.com/soyunomas/nxsh/pkg/types" ha sido eliminada
	"strings"
)

// Node es la interfaz base para todos los nodos del AST.
type Node interface {
	TokenLiteral() string // Devuelve el literal del token asociado al nodo.
	String() string       // Devuelve una representación en string del nodo para depuración.
}

// Statement es una interfaz para nodos de declaración (ej: let x = 5;).
type Statement interface {
	Node
	statementNode()
}

// Expression es una interfaz para nodos de expresión (ej: 5, x, x + 5).
type Expression interface {
	Node
	expressionNode()
}

// Program es el nodo raíz de cualquier AST de nsh.
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}

// LetStatement representa una declaración 'let <ident> = <expresión>;'
type LetStatement struct {
	Token Token      // el token 'let'
	Name  *Identifier
	Value Expression
}

func (ls *LetStatement) statementNode()       {}
func (ls *LetStatement) TokenLiteral() string { return ls.Token.Literal }
func (ls *LetStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ls.TokenLiteral() + " ")
	out.WriteString(ls.Name.String())
	out.WriteString(" = ")
	if ls.Value != nil {
		out.WriteString(ls.Value.String())
	}
	return out.String()
}

// Identifier representa un identificador (nombre de variable, comando, etc).
type Identifier struct {
	Token Token // el token IDENT
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// ExpressionStatement es una declaración que consiste en una única expresión.
// Por ejemplo, `ls -l` es una ExpressionStatement.
type ExpressionStatement struct {
	Token      Token // el primer token de la expresión
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// CommandExpression representa un comando con sus argumentos.
type CommandExpression struct {
	Token Token        // El primer token, que es el nombre del comando.
	Name  Expression   // Será un Identifier.
	Args  []Expression // Argumentos del comando.
}

func (ce *CommandExpression) expressionNode()      {}
func (ce *CommandExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CommandExpression) String() string {
	var out bytes.Buffer
	parts := []string{ce.Name.String()}
	for _, arg := range ce.Args {
		parts = append(parts, arg.String())
	}
	out.WriteString(strings.Join(parts, " "))
	return out.String()
}

// PipelineExpression representa dos comandos conectados por un pipe.
type PipelineExpression struct {
	Token Token      // El token '|'
	Left  Expression
	Right Expression
}

func (pe *PipelineExpression) expressionNode()      {}
func (pe *PipelineExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PipelineExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Left.String())
	out.WriteString(" | ")
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// StringLiteral representa un valor de cadena de texto.
type StringLiteral struct {
	Token Token
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return `"` + sl.Token.Literal + `"` }

// LA FUNCIÓN ToCommand HA SIDO ELIMINADA DE AQUÍ
