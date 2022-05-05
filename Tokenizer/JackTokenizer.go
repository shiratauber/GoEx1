package main

import (
	"bufio"
	"log"
	"os"
)

type Tokenizer struct {
	inputFile string
	File      *os.File
	token     string
}

func New(path string) Tokenizer {
	f, err := os.Open(path)
	Check(err)
	input := Tokenizer{path, f, ""}
	return input
}

func HasMoreTokens(reader *bufio.Reader) bool {
	return reader.Size() > 0
}

func Advance(reader *bufio.Reader, t *Tokenizer) {
	var token string = ""
	for reader.Size() > 0 {
		tav, err := reader.ReadByte()
		Check(err)
		for TabIsDigitsRegex(tav) {
			token += string(tav)

		}

	}

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

func WordIsKeyWord(a string) bool {
	var l []string
	l = append(l, "class", "constructor", "function", "method", "field", "static", "var",
		"int", "char", "boolean", "void", "true", "false", "null", "this", "let", "do", "if",
		"else", "while", "return")

	for _, b := range l {
		if b == a {
			return true
		}
	}
	return false
}

func TabIsSymbol(a byte) bool {
	var l []byte
	l = append(l, '{', '}', '(', ')', '[', ']', '.', ',', ';', '+', '-', '*', '/', '&', '|',
		'<', '>', '=', '~')

	for _, b := range l {
		if b == a {
			return true
		}
	}
	return false
}

func TabIsIdCharsRegex(a byte) bool {
	var l []byte
	l = append(l, 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p',
		'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
		'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9')

	for _, b := range l {
		if b == a {
			return true
		}
	}
	return false
}

func TabIsDigitsRegex(a byte) bool {
	var l []byte
	l = append(l, '0', '1', '2', '3', '4', '5', '6', '7', '8', '9')

	for _, b := range l {
		if b == a {
			return true
		}
	}
	return false
}

func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
