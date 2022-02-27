package parser

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type parser struct {
	inputFile string
	File      *os.File
	Current   string
}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func New(path string) parser {
	f, err := os.Open(path)
	Check(err)
	input := parser{path, f, " "}
	return input
}

func HasMoreCommands(p parser, s *bufio.Scanner) bool {
	return s.Scan()
}

func Advance(s *bufio.Scanner, p parser) {
	p.Current = s.Text()
}

func CommandType(p parser) string {
	var splitt []string = strings.Split(p.Current, " ")

	if WordInArithmetic(splitt[0]) {
		return "C_ARITHMETIC"
	} else if splitt[0] == "push" {
		return "C_PUSH"
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

func Arg1(p parser) string {
	var splitt []string = strings.Split(p.Current, " ")
	if CommandType(p) == "C_ARITHMETIC" {
		return splitt[0]
	} else {
		return splitt[1]
	}

}

func Arg2(p parser) int {
	var splitt []string = strings.Split(p.Current, " ")
	intVar, err := strconv.Atoi(splitt[2])
	Check(err)
	return intVar
}
