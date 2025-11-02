package evaluator

// Environment guarda los identificadores (variables) y sus valores.
type Environment struct {
	store map[string]Object
}

// NewEnvironment crea un nuevo entorno de variables vacío.
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s}
}

// Get recupera un objeto del entorno por su nombre.
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	return obj, ok
}

// Set guarda un objeto en el entorno con un nombre específico.
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}
