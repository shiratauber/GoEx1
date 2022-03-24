package parser

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type Parser struct {
	inputFile  string
	File       *os.File
	Current    string
	LineNumber int
}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func New(path string) Parser {
	f, err := os.Open(path)

	Check(err)
	input := Parser{path, f, " ", 0}

	return input
}

func HasMoreCommands(p *Parser, s *bufio.Scanner) bool {
	return s.Scan()
}

func Advance(s *bufio.Scanner, p *Parser) {
	p.Current = s.Text()
	p.LineNumber += 1
}

func CommandType(p Parser) string {
	var splitt []string = strings.Split(p.Current, " ")

	if WordInArithmetic(splitt[0]) {
		return "C_ARITHMETIC"
	}

	if splitt[0] == "push" {
		return "C_PUSH"
	} else if splitt[0] == "pop" {
		return "C_POP"
	} else if splitt[0] == "label" {
		return "C_LABEL"
	} else if splitt[0] == "goto" {
		return "C_GOTO"
	} else if splitt[0] == "if-goto" {
		return "C_IF"
	} else if splitt[0] == "return" {
		return "C_RETURN"
	} else if splitt[0] == "function" {
		return "C_FUNCTION"
	} else if splitt[0] == "call" {
		return "C_CALL"
	} else {
		return " "
	}

}
func WordInArithmetic(a string) bool {
	var l []string
	l = append(l, "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not")

	for _, b := range l {
		if b == a {
			return true
		}
	}
	return false
}

func Arg1(p Parser) string {
	var splitt []string = strings.Split(p.Current, " ")
	if CommandType(p) == "C_ARITHMETIC" {
		return splitt[0]
	} else {
		return splitt[1]
	}

}

func Arg2(p Parser) int {
	var splitt []string = strings.Split(p.Current, " ")
	intVar, err := strconv.Atoi(splitt[2])
	Check(err)
	return intVar
}
