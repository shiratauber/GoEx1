package main

import (
	"GoEx1/parser"
	"fmt"
	//"GoEx1/CodeWriter"
)

func main() {
	var path string
	fmt.Scanln(&path)
	p := parser.New(path)
	_ = p

}
