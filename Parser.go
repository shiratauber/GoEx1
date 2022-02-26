package GoEx1

import (
	"bufio"
	"log"
	"os"
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

func hasMoreCommands(s *bufio.Scanner) bool {
	return s.Scan()
}

func advance(s *bufio.Scanner, p parser) {
	p.current = s.Text()
}

func commandType() string {
	return " "
}
