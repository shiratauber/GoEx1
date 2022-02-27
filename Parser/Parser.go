package parser

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type parser struct {
	inputFile string
	file      *os.File
	current   string
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func New(path string) parser {
	f, err := os.Open(path)
	check(err)
	input := parser{path, f, " "}
	return input
}

func hasMoreCommands(p parser, s *bufio.Scanner) bool {
	return s.Scan()
}

func advance(s *bufio.Scanner, p parser) {
	p.current = s.Text()
}

func commandType(p parser) string {
	var splitt []string = strings.Split(p.current, " ")

	if wordInArithmetic(splitt[0]) {
		return "C_ARITHMETIC"
	} else {
		return " "
	}
}
func wordInArithmetic(a string) bool {
	var l []string
	l = append(l, "add", "sub", "neg", "eq", "gt", "lt", "and", "or", "not")

	for _, b := range l {
		if b == a {
			return true
		}
	}
	return false
}

func arg1(p parser) string {
	if commandType(p) == "C_ARITHMETIC" {
		var splitt []string = strings.Split(p.current, " ")
		return splitt[0]
	} else {
		return ""
	}

}

func arg2(p parser) int {
	return 0

}
