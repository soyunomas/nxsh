// main.go
package main

import (
	"nxsh/pkg/shell"
)

func main() {
	s := shell.New()
	s.Run()
}
