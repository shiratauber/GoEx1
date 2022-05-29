package jackTokenizer

import (
	"bufio"
	"os"
	"strings"
)

type jackTokenizer struct {
	currentToken     string
	currentTokenType string
	pointer          int
	tokens           []string //  List<String
	tokensString     string
}

func New(inputFile *os.File) jackTokenizer {
	var allCharacters string = ""
	var current string
	var index int
	scan := bufio.NewScanner(inputFile)
	for scan.Scan() {
		current = scan.Text()
		allCharacters += current + "\n"
	}
	mofa := jackTokenizer{"character", "character", 0, nil, ""}

	index = 1
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
			number := character

			for true {
				index = index + 1
				character = allCharacters[index]
				if TabIsDigitsRegex(string(character)) {
					number += character
				} else {
					break
				}
			}

			mofa.tokensString += string(number) + "\n"
		} else if TabIsIdCharsRegex(string(character)) { //   an alphanumeric character(not number, if passed \d)
			word := character

			for true {
				index = index + 1
				character = allCharacters[index]
				if TabIsIdCharsRegex(string(character)) {
					word += character
				} else {
					break
				}
			}

			mofa.tokensString += string(word) + "\n"
		} else if character == '"' {
			//index <- writeStringConstant(allCharacters, index, currentOutputFile)
			stringContent := character
			//   stringContent <- ""
			//  TODO make sure if string should be with or without ""

			for true {
				index = index + 1
				character = allCharacters[index]
				if character == '"' {
					stringContent += character
					index = index + 1 //   jump the closing "
					break
				} else {
					stringContent += character
				}
			}

			mofa.tokensString += string(stringContent) + "\n"
		} else { // ignore white spaces
			index = index + 1
		}
	}
	mofa.pointer = 1
	mofa.currentToken = ""
	mofa.currentTokenType = "NONE"
	mofa.tokens = strings.Split(mofa.tokensString, "\n")

	inputFile.Close()
	return mofa
}

// Do we have more tokens in the input?
func hasMoreTokens(mofa jackTokenizer) bool {
	return mofa.pointer < (len(mofa.tokens) + 1)
}

/* Gets the next token from the input and makes it the current token.
This method should only be called if hasMoreTokens() is true.
Initially there is no current token.*/
func advance(mofa jackTokenizer) {
	if hasMoreTokens(mofa) {
		mofa.currentToken = mofa.tokens[mofa.pointer]
		//   print(paste("DEBUG :", self$currentToken))
		mofa.pointer = mofa.pointer + 1
	} else {
		//print("No more tokens.")
		//quit()
		return
	}

	if WordIsKeyWord(mofa.currentToken) {
		mofa.currentTokenType = "KEYWORD"
	} else if TabIsSymbol(mofa.currentToken) {
		mofa.currentTokenType = "SYMBOL"
	} else if TabIsDigitsRegex(mofa.tokensString) {
		mofa.currentTokenType = "INT_CONST"
	} else if TabIsStringConst(mofa.tokensString) {
		// TODO need to verify if this works
		mofa.currentTokenType = "STRING_CONST"
	} else if TabIsIdCharsRegex(mofa.tokensString) {
		mofa.currentTokenType = "IDENTIFIER"
	} else {
		//print(paste("Unknown token:", self$currentToken))
	}

	// print(paste("DEBUG :", self$currentTokenType))
}

// Returns the type of the current token.
func tokenType(mofa jackTokenizer) string {
	return mofa.currentTokenType
}

// Returns the keyword which is the current token.
// Should be called only when tokeyType() is KEYWORD.
func keyWord(mofa jackTokenizer) string {
	if mofa.currentTokenType == "KEYWORD" {
		return strings.ToUpper(mofa.currentToken)
	} else {
		//   print("Current token is not a keyword!")
		return ""
	}
}

// Returns the character which is the current token.
// Should be called only when tokenType() is SYMBOL.
func symbol(mofa jackTokenizer) string {
	if mofa.currentTokenType == "SYMBOL" {
		//    return(substr(self$currentToken, 1, 1))      ## currentToken[0]
		return mofa.currentToken
	} else {
		//    print("Current token is not a symbol!")
		return ""
	}
}

// Return the identifier which is the current token.
// Should be called only when tokenType() is IDENTIFIER.
func identifier(mofa jackTokenizer) string {
	if mofa.currentTokenType == "IDENTIFIER" {
		return mofa.currentToken
	} else {
		//    print("Current token is not an identifier!")
		return ""
	}
}

// Returns the integer value of the current token.
// Should be called only when tokenType() is INT_CONST.
func intVal(mofa jackTokenizer) string {
	if mofa.currentTokenType == "INT_CONST" {
		return mofa.currentToken
	} else {
		//   print("Current token is not an integer constant!")
		return ""
	}
}

// Returns the string value of the current token without the double quotes.
// Should be called only when tokenType() is STRING_CONST.
func stringVal(mofa jackTokenizer) string {
	if mofa.currentTokenType == "STRING_CONST" {
		size := len(mofa.currentToken)
		return mofa.currentToken[1 : size-1]
	} else {
		//    print("Current token is not a string constant!")
		return ""
	}
}

//  Moves pointer back.
func pointerBack(mofa jackTokenizer) {
	if mofa.pointer > 1 {
		mofa.pointer = mofa.pointer - 1
		mofa.currentToken = mofa.tokens[mofa.pointer]
	}
}

//////////////////////////////////////////////////////////////////////////////
// Returns if current symbol is an op.
func isOp(a string) bool {
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
	l = append(l, '\\', '^', '\n')

	for _, b := range l {
		if string(b) == a {
			return true
		}
	}
	return false
}
