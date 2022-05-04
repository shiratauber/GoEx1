package main

import (
	"bufio"
	"os"
)

func main() {

}

type Tokenizer struct {
	inputFile string
	File      *os.File
}

func New(path string) Tokenizer {

	input := Parser{"", null}

	return input
}

func HasMoreTokens(t *Tokenizer, s *bufio.Scanner) bool {
	return s.Scan()
}
func Advance(s *bufio.Scanner, t *Tokenizer) {
}

func TokenType(t *Tokenizer) string {

	return ""
}

func KeyWord(t *Tokenizer) string {

	return ""
}
func Symbol(t *Tokenizer) string {

	return ""
}

func Identifier(t *Tokenizer) string {

	return ""
}

func IntVal(t *Tokenizer) int {

	return 0
}
func StringVal(t *Tokenizer) string {

	return ""
}
