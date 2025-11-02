package shell

import (
	"encoding/json"
	"fmt"
	"github.com/soyunomas/nxsh/pkg/evaluator"
	"github.com/soyunomas/nxsh/pkg/parser"
	"os"
	"strings"
)

// Constantes para los códigos de color ANSI.
const (
	colorReset = "\033[0m"
	colorCyan  = "\033[36m"
	colorGreen = "\033[32m"
)

type Shell struct {
	// Reemplazamos el antiguo mapa de strings por el nuevo Environment del evaluador.
	environment *evaluator.Environment
	lineReader  LineReader
}

func New() *Shell {
	lr, err := NewLineReader()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error inicializando el lector de línea: %v\n", err)
		os.Exit(1)
	}

	return &Shell{
		// Creamos un único entorno que persistirá durante toda la sesión.
		environment: evaluator.NewEnvironment(),
		lineReader:  lr,
	}
}

func (s *Shell) Run() {
	fmt.Println("Bienvenido a Nexus Shell (nxsh) v1.0.0-rc1.")
	defer s.lineReader.Close()

	for {
		prompt := s.getPrompt()
		line, err := s.lineReader.ReadLine(prompt)
		if err != nil {
			break
		}

		trimmedLine := strings.TrimSpace(line)
		if trimmedLine == "" {
			continue
		}

		// *** LA CORRECCIÓN ESTÁ AQUÍ ***
		// Se compara con la cadena de texto "exit", no con una variable.
		if trimmedLine == "exit" {
			break
		}

		s.lineReader.AddHistory(trimmedLine)
		s.eval(trimmedLine)
	}
}

func (s *Shell) getPrompt() string {
	wd, err := os.Getwd()
	if err != nil {
		return "nxsh > "
	}
	home, _ := os.UserHomeDir()
	if home != "" {
		wd = strings.Replace(wd, home, "~", 1)
	}
	return fmt.Sprintf("%s%s%s %s%s%s ", colorCyan, wd, colorReset, colorGreen, "nxsh >", colorReset)
}

func (s *Shell) eval(line string) {
	l := parser.NewLexer(line)
	p := parser.NewParser(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Fprintln(os.Stderr, "Error de parsing:", msg)
		}
		return
	}

	// El evaluador ahora necesita el entorno del shell para operar.
	evaluated := evaluator.Eval(program, s.environment)

	// El evaluador devuelve NULL para 'let', no debemos imprimir nada en ese caso.
	if evaluated != nil && evaluated.Type() != evaluator.NULL_OBJ {
		switch obj := evaluated.(type) {
		case *evaluator.Error:
			fmt.Fprintln(os.Stderr, obj.Inspect())
		case *evaluator.String:
			fmt.Print(obj.Inspect())
			if !strings.HasSuffix(obj.Inspect(), "\n") {
				fmt.Println()
			}
		case *evaluator.Json:
			encoder := json.NewEncoder(os.Stdout)
			encoder.SetIndent("", "  ")
			if err := encoder.Encode(obj.Value); err != nil {
				fmt.Fprintln(os.Stderr, "Error al formatear JSON:", err)
			}
		default:
			fmt.Println(obj.Inspect())
		}
	}
}
