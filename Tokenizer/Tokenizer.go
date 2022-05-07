package Tokenizer

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Tokenizer struct {
	//InputFile string
	InputFile     *os.File
	OutputFile    *os.File
	Token         string
	AllCharacters string
}

func New(path string) Tokenizer {
	inputFile, err := os.Open(path)
	Check(err)

	var split []string = strings.Split(path, "\\")
	var last string = split[len(split)-1] //the name of the output file
	var splitt []string = strings.Split(last, ".")
	last = splitt[0]
	//create the output file
	outputFile, err := os.Create(last + "TT.xml")
	Check(err)
	//open the output file
	myFile, err := os.OpenFile(outputFile.Name(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	Check(err)

	mofa := Tokenizer{inputFile, myFile, "", ""}
	return mofa
}

func Advance(scan *bufio.Scanner, t Tokenizer) {
	var current string
	var index int
	for scan.Scan() {
		current = scan.Text()
		t.AllCharacters += current + "\n"
	}

	var character uint8
	index = 0
	var sizeOfFile int = len(t.AllCharacters)
	for index < sizeOfFile {
		character = t.AllCharacters[index]
		if character == '/' {
			nextCharacter := t.AllCharacters[index+1]
			if nextCharacter == '/' {
				index = index + 1
				index = handleSingleLineComment(t.AllCharacters, index)
			} else if nextCharacter == '*' {
				index = index + 1
				index = handleMultipleLinesComment(t.AllCharacters, index)
			} else {
				index = writeSymbol(character, index, t.OutputFile)
			}
		} else if TabIsSymbol(character) {
			index = writeSymbol(character, index, t.OutputFile)
		} else if TabIsDigitsRegex(character) {
			index = writeIntegerConstant(t.AllCharacters, index, t.OutputFile)
		} else if TabIsIdCharsRegex(character) {
			index = writeIdentifierOrKeyword(t.AllCharacters, index, t.OutputFile)
		} else if character == '"' {
			index = writeStringConstant(t.AllCharacters, index, t.OutputFile)
		} else { // ignore white spaces
			index = index + 1
		}
	}

}

func writeStringConstant(characters string, index int, file *os.File) int {
	stringContent := ""

	for true {
		index = index + 1
		character := characters[index]
		if character == '"' {
			index = index + 1 //jump the closing "
			break
		}
		stringContent += string(character)
	}

	writeToXML("stringConstant", stringContent, file)
	return index
}

func writeIdentifierOrKeyword(characters string, index int, file *os.File) int {
	word := string(characters[index])

	for true {
		index = index + 1
		character := characters[index]
		if !TabIsIdCharsRegex(character) {
			break
		}
		word += string(character)
	}

	if WordIsKeyWord(word) {
		writeToXML("keyword", word, file)
	} else {
		writeToXML("identifier", word, file)
	}

	return index
}

func writeIntegerConstant(characters string, index int, file *os.File) int {
	number := string(characters[index])

	for true {
		index = index + 1
		character := characters[index]
		if !TabIsDigitsRegex(character) {
			break
		}
		number += string(character)
	}
	writeToXML("integerConstant", number, file)
	return index
}

func writeSymbol(character uint8, index int, file *os.File) int {
	index = index + 1
	var stringCharacter string
	stringCharacter = string(character)
	if character == '<' {
		stringCharacter = "&lt;"
	}
	if character == '>' {
		stringCharacter = "&gt;"
	}
	if character == '&' {
		stringCharacter = "&amp;"
	}

	writeToXML("symbol", stringCharacter, file)
	return index
}

func writeToXML(tag string, content string, file *os.File) {
	// print(paste(tag, content, sep=" : "))
	var s string = "<" + tag + "> " + content + " </" + tag + ">" + "\n"
	if _, err := file.WriteString(s); err != nil {
		panic(err)
	}
}

func handleMultipleLinesComment(characters string, index int) int {
	for true {
		index = index + 1
		character := characters[index]
		if character == '*' {
			nextCharacter := characters[index+1]
			if nextCharacter == '/' {
				index = index + 2 // jump the closing */
				break
			}
		}
	}

	return index
}

func handleSingleLineComment(characters string, index int) int {
	for true {
		index = index + 1
		character := characters[index]
		if character == '\n' {
			index = index + 1 //jump the end line
			break
		}
	}

	return index
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

func Close(t Tokenizer) {
	t.InputFile.Close()
	t.OutputFile.Close()
}
