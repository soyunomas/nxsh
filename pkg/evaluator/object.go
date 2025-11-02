package evaluator

import (
	"encoding/json"
	"fmt"
)

// ObjectType es el tipo de un objeto en nsh.
type ObjectType string

const (
	STRING_OBJ  ObjectType = "STRING"
	JSON_OBJ    ObjectType = "JSON"
	NULL_OBJ    ObjectType = "NULL"
	ERROR_OBJ   ObjectType = "ERROR"
	BUILTIN_OBJ ObjectType = "BUILTIN"
)

// Object es la interfaz que todo tipo de dato en nsh debe implementar.
type Object interface {
	Type() ObjectType
	Inspect() string
}

// BuiltinFunction es el tipo de las funciones internas de nsh.
type BuiltinFunction func(input Object, args ...Object) Object

// Builtin representa una función interna.
type Builtin struct {
	Fn BuiltinFunction
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }

// String representa un valor de cadena.
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }

// Json representa datos JSON parseados.
type Json struct {
	Value interface{}
}

func (j *Json) Type() ObjectType { return JSON_OBJ }
func (j *Json) Inspect() string {
	// Usamos MarshalIndent para que el JSON tenga saltos de línea
	// y pueda ser procesado por herramientas como grep.
	b, err := json.MarshalIndent(j.Value, "", "  ") // "" para prefijo, "  " para indentación
	if err != nil {
		return fmt.Sprintf("Error marshaling JSON: %v", err)
	}
	return string(b)
}

// Null representa la ausencia de valor.
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// Error representa un error que ocurrió durante la evaluación.
type Error struct {
	Message string
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "Error: " + e.Message }
