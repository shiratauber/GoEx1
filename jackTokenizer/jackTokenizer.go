package jackTokenizer

import (
	"GoEx1/Analyzer"
	"bufio"
	"os"
	"strconv"
	"strings"
)

type JackTokenizer struct {
	CurrentToken     string
	currentTokenType string
	pointer          int
	tokens           []string //  List<String>
	tokensString     string
}

func New(inputFile *os.File) JackTokenizer {
	var allCharacters string = ""
	var current string
	var index int
	scan := bufio.NewScanner(inputFile)
	for scan.Scan() {
		current = scan.Text()
		allCharacters += current + "\n"
	}
	mofa := JackTokenizer{"character", "character", 0, nil, ""}

	index = 0
	var totalNumberOfCharacters = len(allCharacters)

	for index < totalNumberOfCharacters {
		character := allCharacters[index]

		if character == '/' {
			nextCharacter := allCharacters[index+1]
			if nextCharacter == '/' { //   // comments
				index = index + 1
				for true {
					index = index + 1
					character = allCharacters[index]
					if character == '\n' {
						index = index + 1 //   jump the endline
						break

					}
				}
			} else if nextCharacter == '*' { //  /* comments */
				index = index + 1
				for true {
					index = index + 1
					character = allCharacters[index]
					if character == '*' {
						nextCharacter = allCharacters[index+1]
						if nextCharacter == '/' {
							index = index + 2 //   jump the closing */
							break
						}
					}
				}
			} else { //   just the / symbol
				mofa.tokensString += string(character) + "\n"
				index = index + 1
			}
		} else if TabIsSymbol(string(character)) {
			mofa.tokensString += string(character) + "\n"
			index = index + 1
		} else if TabIsDigitsRegex(string(character)) { //   a digit
			var number string = string(character)

			for true {
				index = index + 1
				character = allCharacters[index]
				if TabIsDigitsRegex(string(character)) {
					number += string(character)
				} else {
					break
				}
			}

			mofa.tokensString += number + "\n"
		} else if TabIsIdCharsRegex(string(character)) { //   an alphanumeric character(not number, if passed \d)
			var word string = string(character)

			for true {
				index = index + 1
				character = allCharacters[index]
				if TabIsIdCharsRegex(string(character)) {
					word += string(character)
				} else {
					break
				}
			}

			mofa.tokensString += word + "\n"
		} else if character == '"' {
			//index <- writeStringConstant(allCharacters, index, currentOutputFile)
			var stringContent string = string(character)
			//   stringContent <- ""
			//  TODO make sure if string should be with or without ""

			for true {
				index = index + 1
				character = allCharacters[index]
				if character == '"' {
					stringContent += string(character)
					index = index + 1 //   jump the closing "
					break
				} else {
					stringContent += string(character)
				}
			}

			mofa.tokensString += stringContent + "\n"
		} else { // ignore white spaces
			index = index + 1
		}
	}
	mofa.pointer = 0
	mofa.CurrentToken = ""
	mofa.currentTokenType = "NONE"
	mofa.tokens = strings.Split(mofa.tokensString, "\n")

	inputFile.Close()
	return mofa
}

// Do we have more tokens in the input?
func HasMoreTokens(mofa *JackTokenizer) bool {
	return mofa.pointer < (len(mofa.tokens) + 1)
}

// Gets the next token from the input and makes it the current token.
// This method should only be called if hasMoreTokens() is true.
// Initially there is no current token.
func Advance(mofa *JackTokenizer) {
	if HasMoreTokens(mofa) {
		mofa.CurrentToken = mofa.tokens[mofa.pointer]
		//   print(paste("DEBUG :", self$currentToken))
		mofa.pointer = mofa.pointer + 1
	} else {
		//print("No more tokens.")
		//quit()
		return
	}

	if WordIsKeyWord(mofa.CurrentToken) {
		mofa.currentTokenType = "KEYWORD"
	} else if TabIsSymbol(mofa.CurrentToken) {
		mofa.currentTokenType = "SYMBOL"
	} else if TabIsDigitsRegex(string(mofa.CurrentToken[0])) {
		mofa.currentTokenType = "INT_CONST"
	} else if TabIsStringConst(string(mofa.CurrentToken[0])) {
		// TODO need to verify if this works
		mofa.currentTokenType = "STRING_CONST"
	} else if TabIsIdCharsRegex(string(mofa.CurrentToken[0])) {
		mofa.currentTokenType = "IDENTIFIER"
	} else {
		print("Unknown token:" + mofa.CurrentToken)
	}

	// print(paste("DEBUG :", self$currentTokenType))
}

// Returns the type of the current token.
func TokenType(mofa JackTokenizer) string {
	return mofa.currentTokenType
}

// Returns the keyword which is the current token.
// Should be called only when tokeyType() is KEYWORD.
func KeyWord(mofa JackTokenizer) string {
	if mofa.currentTokenType == "KEYWORD" {
		return strings.ToUpper(mofa.CurrentToken)
	} else {
		//   print("Current token is not a keyword!")
		return ""
	}
}

// Returns the character which is the current token.
// Should be called only when TokenType() is SYMBOL.
func Symbol(mofa JackTokenizer) string {
	if mofa.currentTokenType == "SYMBOL" {
		//    return(substr(self$currentToken, 1, 1))      ## currentToken[0]
		return strings.ToUpper(mofa.CurrentToken)
	} else {
		//    print("Current token is not a symbol!")
		return ""
	}
}

// Return the Identifier which is the current token.
// Should be called only when TokenType() is IDENTIFIER.
func Identifier(mofa JackTokenizer) string {
	if mofa.currentTokenType == "IDENTIFIER" {
		return mofa.CurrentToken
	} else {
		//    print("Current token is not an identifier!")
		return ""
	}
}

// Returns the integer value of the current token.
// Should be called only when TokenType() is INT_CONST.
func IntVal(mofa JackTokenizer) int {
	if mofa.currentTokenType == "INT_CONST" {
		var intVar, err = strconv.Atoi(mofa.CurrentToken)
		Analyzer.Check(err)
		return intVar
	} else {
		//   print("Current token is not an integer constant!")
		return -1
	}
}

// Returns the string value of the current token without the double quotes.
// Should be called only when TokenType() is STRING_CONST.
func StringVal(mofa JackTokenizer) string {
	if mofa.currentTokenType == "STRING_CONST" {
		size := len(mofa.CurrentToken)
		return mofa.CurrentToken[1 : size-1]
	} else {
		//    print("Current token is not a string constant!")
		return ""
	}
}

//  Moves pointer back.
func PointerBack(mofa *JackTokenizer) {
	if mofa.pointer > 1 {
		mofa.pointer = mofa.pointer - 1
		mofa.CurrentToken = mofa.tokens[mofa.pointer]
	}
}

//////////////////////////////////////////////////////////////////////////////
// Returns if current Symbol is an op.
func IsOp(a string) bool {
	//return (symbol in {'+', '-', '*', '/', '&', '|', '<', '>', '='})
	var l []byte
	l = append(l, '+', '-', '*', '/', '&', '|', '<', '>', '=')

	for _, b := range l {
		if string(b) == a {
			return true
		}
	}
	return false
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

func TabIsSymbol(a string) bool {
	var l []byte
	l = append(l, '{', '}', '(', ')', '[', ']', '.', ',', ';', '+', '-', '*', '/', '&', '|',
		'<', '>', '=', '~')

	for _, b := range l {
		if string(b) == a {
			return true
		}
	}
	return false
}

func TabIsIdCharsRegex(a string) bool {
	var l []byte
	l = append(l, 'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j', 'k', 'l', 'm', 'n', 'o', 'p',
		'q', 'r', 's', 't', 'u', 'v', 'w', 'x', 'y', 'z', 'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J',
		'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'R', 'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z',
		'0', '1', '2', '3', '4', '5', '6', '7', '8', '9')

	for _, b := range l {
		if string(b) == a {
			return true
		}
	}
	return false
}

func TabIsDigitsRegex(a string) bool {
	var l []byte
	l = append(l, '0', '1', '2', '3', '4', '5', '6', '7', '8', '9')

	for _, b := range l {
		if string(b) == a {
			return true
		}
	}
	return false
}

/////////////////////////////////////////////////////////////////////////////////
func TabIsStringConst(a string) bool {
	var l []byte
	l = append(l, '"')

	for _, b := range l {
		if string(b) == a {
			return true
		}
	}
	return false
}
