package main

import (
	"os"

	"github.com/sam8helloworld/uwscgo/repl"
)

func main() {
	repl.Start(os.Stdin, os.Stdout)
}
