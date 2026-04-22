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
	data, err := cli.Parse[CLIStruct](os.Args)
	if err != nil {
		panic(err)
	}
	fmt.Printf("%#v\n", data)
}
