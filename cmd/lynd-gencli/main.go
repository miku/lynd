package main

import (
	"flag"
	"fmt"
	"go/parser"
	"go/token"
	"log"
	"os"
)

func main() {
	dir := flag.String("d", ".", "root dir for finding files")

	flag.Parse()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, *dir, func(os.FileInfo) bool { return true }, 0)
	if err != nil {
		log.Fatal(err)
	}
	for k, v := range pkgs {
		fmt.Printf("%s -> %+v\n", k, v)
	}
}
