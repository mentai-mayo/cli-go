package main

import (
	"fmt"
	"os"

	"github.com/mentai-mayo/cli-go"
)

type CLIStruct struct {
	Name    string `pos:"1"`
	Version string `short:"v" long:"version"`
}

func main() {
	fmt.Println("Hello, world!")
	cli.Parse[CLIStruct](os.Args)
}
