package shell

import (
	"io"
	"os"

	"github.com/chzyer/readline"
)

// LineReader es una interfaz que abstrae la lectura de líneas interactivas.
type LineReader interface {
	ReadLine(prompt string) (string, error)
	Close() error
	AddHistory(line string)
}

// nshCompleter implementa la interfaz de autocompletado de readline.
type nshCompleter struct{}

func (nc *nshCompleter) Do(line []rune, pos int) ([][]rune, int) {
	// Por ahora, no implementamos autocompletado.
	// Esto es el placeholder para el futuro.
	return nil, 0
}

// NshReadline es nuestra implementación de LineReader usando la librería readline.
type NshReadline struct {
	instance *readline.Instance
}

// NewLineReader crea y configura una nueva instancia de NshReadline.
func NewLineReader() (*NshReadline, error) {
	completer := &nshCompleter{}

	// Configurar el historial de comandos
	home, _ := os.UserHomeDir()
	historyFile := ""
	if home != "" {
		historyFile = home + "/.nxsh_history"
	}

	config := &readline.Config{
		Prompt:          "> ", // El prompt se establecerá dinámicamente
		HistoryFile:     historyFile,
		AutoComplete:    completer,
		InterruptPrompt: "^C",
		EOFPrompt:       "exit",
	}

	rl, err := readline.NewEx(config)
	if err != nil {
		return nil, err
	}

	return &NshReadline{instance: rl}, nil
}

// ReadLine lee una línea de la entrada estándar con el prompt dado.
func (nr *NshReadline) ReadLine(prompt string) (string, error) {
	nr.instance.SetPrompt(prompt)
	line, err := nr.instance.Readline()

	if err == readline.ErrInterrupt {
		// Si el usuario pulsa Ctrl-C en una línea vacía, readline devuelve ErrInterrupt.
		// Si la línea no está vacía, la limpia. Queremos tratarlo como una línea vacía.
		return "", nil
	} else if err == io.EOF {
		// Ctrl-D
		return "exit", nil
	}

	return line, err
}

// Close cierra la instancia de readline.
func (nr *NshReadline) Close() error {
	return nr.instance.Close()
}

// AddHistory añade una línea al historial de comandos.
func (nr *NshReadline) AddHistory(line string) {
	nr.instance.SaveHistory(line)
}
