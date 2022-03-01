package main

/*
Shira Tauber 213936271
Shvut Lazare 213195977
exercise 1
group 43
*/

import (
	"GoEx1/CodeWriter"
	"GoEx1/Parser"
	"bufio"
	"fmt"
	"strconv"
)

func main() {
	var path string
	fmt.Scanln(&path)
	pars := parser.New(path)
	code := CodeWriter.New(path)
	scanner := bufio.NewScanner(pars.File)
	for true {
		if parser.HasMoreCommands(&pars, scanner) {
			parser.Advance(scanner, &pars)
			if parser.CommandType(pars) == "C_ARITHMETIC" {
				CodeWriter.WriteArithmetic(pars.Current, code, strconv.Itoa(pars.LineNumber))
			} else if parser.CommandType(pars) == "C_PUSH" || parser.CommandType(pars) == "C_POP" {
				CodeWriter.WritePushPop(parser.CommandType(pars), parser.Arg1(pars), parser.Arg2(pars), code)
			}
		} else {
			break
		}
	}
	CodeWriter.Close(code)
	//_ = p

}
