package evaluator

import (
	"bytes"
	"encoding/json"
	"fmt"
	"nxsh/pkg/parser"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

var (
	NULL = &Null{}
)

var builtins = map[string]*Builtin{
	"cd":     {Fn: builtinCd},
	"get":    {Fn: builtinGet},
	"where":  {Fn: builtinWhere},
	"select": {Fn: builtinSelect}, // <-- REGISTRAMOS SELECT
}

// builtinSelect implementa el comando 'select' para proyectar campos de objetos JSON.
func builtinSelect(input Object, args ...Object) Object {
	if input == nil {
		return newError("select: requiere una entrada de un pipeline")
	}
	jsonInput, ok := input.(*Json)
	if !ok {
		return newError("select: la entrada debe ser de tipo JSON, se obtuvo %s", input.Type())
	}
	if len(args) < 1 {
		return newError("uso: select <.campo1> <.campo2> ...")
	}

	// Pre-parsear todas las rutas de los argumentos
	var paths [][]string
	for _, arg := range args {
		pathArg, ok := arg.(*String)
		if !ok {
			return newError("select: todos los argumentos deben ser cadenas de ruta")
		}
		pathStr := strings.TrimPrefix(pathArg.Value, ".")
		paths = append(paths, strings.Split(pathStr, "."))
	}

	// Función auxiliar para procesar un único objeto
	processObject := func(itemData map[string]interface{}) map[string]interface{} {
		newObj := make(map[string]interface{})
		for _, path := range paths {
			if value, found := accessField(itemData, path); found {
				newKey := path[len(path)-1] // La nueva clave es la última parte de la ruta
				newObj[newKey] = value
			}
		}
		return newObj
	}

	switch data := jsonInput.Value.(type) {
	case map[string]interface{}: // Entrada es un único objeto
		return &Json{Value: processObject(data)}
	case []interface{}: // Entrada es un array de objetos
		var results []interface{}
		for _, item := range data {
			if itemMap, ok := item.(map[string]interface{}); ok {
				newObj := processObject(itemMap)
				if len(newObj) > 0 {
					results = append(results, newObj)
				}
			}
		}
		return &Json{Value: results}
	default:
		return newError("select: solo puede operar sobre objetos o arrays de objetos JSON")
	}
}

// builtinWhere implementa el comando 'where' para filtrar arrays de objetos.
func builtinWhere(input Object, args ...Object) Object {
	if input == nil {
		return newError("where: requiere una entrada de un pipeline")
	}
	jsonInput, ok := input.(*Json)
	if !ok {
		return newError("where: la entrada debe ser de tipo JSON, se obtuvo %s", input.Type())
	}
	items, ok := jsonInput.Value.([]interface{})
	if !ok {
		if itemMap, isMap := jsonInput.Value.(map[string]interface{}); isMap {
			items = []interface{}{itemMap}
		} else {
			return newError("where: solo puede filtrar arrays de objetos")
		}
	}
	if len(args) != 3 {
		return newError("uso: where <.campo> <operador> <valor>")
	}
	pathArg, okPath := args[0].(*String)
	opArg, okOp := args[1].(*String)
	valArg, okVal := args[2].(*String)
	if !okPath || !okOp || !okVal {
		return newError("where: los argumentos deben ser cadenas")
	}
	pathStr := strings.TrimPrefix(pathArg.Value, ".")
	path := strings.Split(pathStr, ".")
	op := opArg.Value
	valueStr := valArg.Value
	var results []interface{}
	for _, item := range items {
		itemValue, found := accessField(item, path)
		if !found {
			continue
		}
		match, err := evaluateCondition(itemValue, op, valueStr)
		if err != nil {
			return newError("error en 'where': %v", err)
		}
		if match {
			results = append(results, item)
		}
	}
	return &Json{Value: results}
}

// evaluateCondition compara un valor de JSON (lhs) con un string (rhsStr).
func evaluateCondition(lhs interface{}, op string, rhsStr string) (bool, error) {
	if lhsFloat, ok := lhs.(float64); ok {
		if rhsFloat, err := strconv.ParseFloat(rhsStr, 64); err == nil {
			switch op {
			case "==": return lhsFloat == rhsFloat, nil
			case "!=": return lhsFloat != rhsFloat, nil
			case ">": return lhsFloat > rhsFloat, nil
			case "<": return lhsFloat < rhsFloat, nil
			case ">=": return lhsFloat >= rhsFloat, nil
			case "<=": return lhsFloat <= rhsFloat, nil
			default: return false, fmt.Errorf("operador numérico desconocido: %s", op)
			}
		}
	}
	if lhsBool, ok := lhs.(bool); ok {
		if rhsBool, err := strconv.ParseBool(rhsStr); err == nil {
			switch op {
			case "==": return lhsBool == rhsBool, nil
			case "!=": return lhsBool != rhsBool, nil
			default: return false, fmt.Errorf("operador booleano no válido: %s", op)
			}
		}
	}
	lhsStr := fmt.Sprintf("%v", lhs)
	rhsStrUnquoted, err := strconv.Unquote(`"` + rhsStr + `"`)
	if err != nil {
		rhsStrUnquoted = rhsStr
	}
	switch op {
	case "==": return lhsStr == rhsStrUnquoted, nil
	case "!=": return lhsStr != rhsStrUnquoted, nil
	default: return false, fmt.Errorf("operador no soportado para el tipo de dato: %s", op)
	}
}

// builtinGet implementa el comando 'get' para extraer datos de objetos JSON.
func builtinGet(input Object, args ...Object) Object {
	if input == nil {
		return newError("get: requiere una entrada de un pipeline")
	}
	jsonInput, ok := input.(*Json)
	if !ok {
		return newError("get: la entrada debe ser de tipo JSON, se obtuvo %s", input.Type())
	}
	if len(args) != 1 {
		return newError("uso: get <.campo.anidado>")
	}
	pathArg, ok := args[0].(*String)
	if !ok {
		return newError("get: el argumento de ruta debe ser una cadena, se obtuvo %s", args[0].Type())
	}
	pathStr := strings.TrimPrefix(pathArg.Value, ".")
	path := strings.Split(pathStr, ".")
	switch data := jsonInput.Value.(type) {
	case map[string]interface{}:
		result, found := accessField(data, path)
		if !found { return NULL }
		return nativeToNshObject(result)
	case []interface{}:
		var results []interface{}
		for _, item := range data {
			if result, found := accessField(item, path); found {
				results = append(results, result)
			}
		}
		return &Json{Value: results}
	default:
		return newError("get: la entrada no es un objeto o array JSON válido")
	}
}

// accessField es una función auxiliar para navegar recursivamente en datos JSON parseados.
func accessField(data interface{}, path []string) (interface{}, bool) {
	if len(path) == 0 { return data, true }
	currentKey := path[0]
	remainingPath := path[1:]
	obj, ok := data.(map[string]interface{})
	if !ok { return nil, false }
	value, found := obj[currentKey]
	if !found { return nil, false }
	return accessField(value, remainingPath)
}

// nativeToNshObject convierte un valor nativo de Go (de JSON) a un objeto de nsh.
func nativeToNshObject(v interface{}) Object {
	if v == nil { return NULL }
	if _, isMap := v.(map[string]interface{}); isMap { return &Json{Value: v} }
	if _, isSlice := v.([]interface{}); isSlice { return &Json{Value: v} }
	return &String{Value: fmt.Sprintf("%v", v)}
}

func builtinCd(_ Object, args ...Object) Object {
	if len(args) > 1 { return newError("cd: demasiados argumentos") }
	var path string
	if len(args) == 0 {
		home, err := os.UserHomeDir()
		if err != nil { return newError("cd: no se pudo encontrar el directorio home: %v", err) }
		path = home
	} else {
		str, ok := args[0].(*String)
		if !ok { return newError("cd: el argumento debe ser una cadena, se obtuvo %s", args[0].Type()) }
		path = str.Value
	}
	if err := os.Chdir(path); err != nil { return newError("cd: %v", err) }
	return NULL
}

func Eval(node parser.Node, env *Environment) Object {
	switch node := node.(type) {
	case *parser.Program: return evalProgram(node, env)
	case *parser.ExpressionStatement: return Eval(node.Expression, env)
	case *parser.LetStatement:
		val := Eval(node.Value, env)
		if isError(val) { return val }
		env.Set(node.Name.Value, val)
		return NULL
	case *parser.Identifier: return evalIdentifier(node, env)
	case *parser.CommandExpression: return evalCommandExpression(node, env, nil)
	case *parser.StringLiteral: return &String{Value: node.Value}
	case *parser.PipelineExpression:
		leftResult := Eval(node.Left, env)
		if isError(leftResult) { return leftResult }
		return evalPipelineChain(node.Right, env, leftResult)
	}
	return newError("tipo de nodo no soportado: %T", node)
}

func evalPipelineChain(node parser.Expression, env *Environment, input Object) Object {
	switch node := node.(type) {
	case *parser.CommandExpression: return evalCommandExpression(node, env, input)
	case *parser.PipelineExpression:
		intermediateResult := evalPipelineChain(node.Left, env, input)
		if isError(intermediateResult) { return intermediateResult }
		return evalPipelineChain(node.Right, env, intermediateResult)
	default: return newError("lado derecho del pipe inválido: se esperaba un comando")
	}
}

func evalProgram(program *parser.Program, env *Environment) Object {
	var result Object
	for _, statement := range program.Statements {
		result = Eval(statement, env)
		if err, ok := result.(*Error); ok { return err }
	}
	return result
}

func evalIdentifier(node *parser.Identifier, env *Environment) Object {
	if val, ok := env.Get(node.Value); ok { return val }
	if builtin, ok := builtins[node.Value]; ok { return builtin }
	return &String{Value: node.Value}
}

func evalCommandExpression(cmdExpr *parser.CommandExpression, env *Environment, input Object) Object {
	if ident, ok := cmdExpr.Name.(*parser.Identifier); ok {
		if val, exists := env.Get(ident.Value); exists {
			if len(cmdExpr.Args) > 0 {
				return newError("la variable '%s' no es un comando y no acepta argumentos", ident.Value)
			}
			return val
		}
	}
	nameObj := Eval(cmdExpr.Name, env)
	if isError(nameObj) { return nameObj }
	var args []Object
	for _, argExpr := range cmdExpr.Args {
		evaluatedArg := Eval(argExpr, env)
		if isError(evaluatedArg) { return evaluatedArg }
		args = append(args, evaluatedArg)
	}
	if builtin, ok := nameObj.(*Builtin); ok { return builtin.Fn(input, args...) }
	cmdName := nameObj.Inspect()
	var argStrings []string
	for _, arg := range args {
		argStrings = append(argStrings, arg.Inspect())
	}
	cmd := exec.Command(cmdName, argStrings...)
	if input != nil {
		cmd.Stdin = bytes.NewReader([]byte(input.Inspect()))
	} else {
		cmd.Stdin = os.Stdin
	}
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return newError("error ejecutando '%s': %v", cmdName, err)
	}
	output := out.Bytes()
	var jsonData interface{}
	if err := json.Unmarshal(output, &jsonData); err == nil {
		return &Json{Value: jsonData}
	}
	return &String{Value: string(output)}
}

func newError(format string, a ...interface{}) *Error { return &Error{Message: fmt.Sprintf(format, a...)} }
func isError(obj Object) bool {
	if obj != nil { return obj.Type() == ERROR_OBJ }
	return false
}
