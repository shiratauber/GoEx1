package Tokenizer

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Tokenizer struct {
	//inputFile string
	inputFile  *os.File
	outputFile *os.File
	token      string
}

func New(path string) Tokenizer {
	inputfile, err := os.Open(path)
	Check(err)

	var splitt []string = strings.Split(path, "\\")
	var last string = splitt[len(splitt)-1] //the name of the output file

	//create the output file
	outputfile, err := os.Create(last + "TT.xml")
	Check(err)
	//open the output file
	myFile, err := os.OpenFile(outputfile.Name(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	Check(err)

	mofa := Tokenizer{inputfile, myFile, ""}
	return mofa
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
