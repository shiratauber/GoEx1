package Analyzer

import (
	"bufio"
	"log"
	"os"
	"strings"
)

type Analyzer struct {
	InputFile     *os.File
	OutputFile    *os.File
	AllCharacters string
}

func New(path string) Analyzer {
	inputFile, err := os.Open(path)
	Check(err)

	var split []string = strings.Split(path, "\\")
	var last string = split[len(split)-1] //the name of the output file
	var splitt []string = strings.Split(last, ".")
	last = splitt[0]
	//create the output file
	outputFile, err := os.Create(last + "New.xml")
	Check(err)
	//open the output file
	myFile, err := os.OpenFile(outputFile.Name(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	Check(err)

	mofa := Analyzer{inputFile, myFile, ""}
	return mofa

}

func WriteClass(scan *bufio.Scanner, a Analyzer) {
	var current string = ""
	for scan.Scan() {
		current += scan.Text() + "\n"
	}
	lines := strings.Split(current, "\n")
	index := 2

	tagLevel := 0
	writeToXML("class", tagLevel, a.OutputFile)
	tagLevel = tagLevel + 1

	index = writeLinesToXML(lines, index, 3, tagLevel, a.OutputFile) // class className {

	for true {
		words := strings.Split(lines[index], " ")[1]
		word := string(words[2]) // word between the tags

		if word == "static" || word == "field" {
			index = writeClassVarDec(lines, index, tagLevel, a.OutputFile) // classVarDec*
		} else if word == "constructor" || word == "function" || word == "method" {
			index = writeSubroutineDec(lines, index, tagLevel, a.OutputFile) // subroutineDec*
		} else {
			break
		}
	}

	index = writeLinesToXML(lines, index, 1, tagLevel, a.OutputFile) // }
	tagLevel = tagLevel - 1
	writeToXML("/class", tagLevel, a.OutputFile)

}
func writeSubroutineDec(lines []string, index int, tagLevel int, output *os.File) int {
	writeToXML("subroutineDec", tagLevel, output)
	tagLevel = tagLevel + 1

	index = writeLinesToXML(lines, index, 4, tagLevel, output)  // write 4 starting lines
	index = writeParameterList(lines, index, tagLevel, output)  // parameterList
	index = writeLinesToXML(lines, index, 1, tagLevel, output)  //)
	index = writeSubroutineBody(lines, index, tagLevel, output) //subroutineBody

	tagLevel = tagLevel - 1
	writeToXML("/subroutineDec", tagLevel, output)

	return index

}
func writeSubroutineBody(lines []string, index int, tagLevel int, output *os.File) int {
	writeToXML("subroutineBody", tagLevel, output)
	tagLevel = tagLevel + 1

	index = writeLinesToXML(lines, index, 1, tagLevel, output) // {

	for true {
		var words = strings.Split(string(lines[index]), " ")[1] ////////////////////////////////////////////
		var word = words[2]                                     //word between the tags

		if string(word) == "var" {
			index = writeVarDec(lines, index, tagLevel, output) // varDec*
		} else {
			break
		}
	}

	index = writeStatements(lines, index, tagLevel, output) // statements
	index = writeLinesToXML(lines, index, 1, tagLevel, output)

	tagLevel = tagLevel - 1
	writeToXML("/subroutineBody", tagLevel, output)

	return index
}
func writeVarDec(lines []string, index int, tagLevel int, output *os.File) int {
	writeToXML("varDec", tagLevel, output)
	tagLevel = tagLevel + 1

	index = writeLinesToXML(lines, index, 3, tagLevel, output) // var type varName

	for true {
		words := strings.Split(lines[index], " ")[1]
		word := words[2] // word between the tags
		if word == ',' {
			index = writeLinesToXML(lines, index, 2, tagLevel, output) // (, varName)*
		} else {
			break
		}
	}

	index = writeLinesToXML(lines, index, 1, tagLevel, output) // ;
	tagLevel = tagLevel - 1
	writeToXML("/varDec", tagLevel, output)

	return index
}
func writeParameterList(lines []string, index int, tagLevel int, output *os.File) int {
	writeToXML("parameterList", tagLevel, output)
	tagLevel = tagLevel + 1

	words := strings.Split(lines[index], " ")[1]
	word := words[2] // word between the tags

	if word != ')' { // not empty list, ?
		index = writeLinesToXML(lines, index, 2, tagLevel, output) // type varName

		for true {
			words := strings.Split(lines[index], " ")[1]
			word := words[2] // word between the tags
			if word == ',' {
				index = writeLinesToXML(lines, index, 3, tagLevel, output) // (, type varName)*
			} else {
				break
			}
		}
	}

	tagLevel = tagLevel - 1
	writeToXML("/parameterList", tagLevel, output)

	return index
}
func writeClassVarDec(lines []string, index int, tagLevel int, output *os.File) int {
	writeToXML("classVarDec", tagLevel, output)
	tagLevel = tagLevel + 1

	index = writeLinesToXML(lines, index, 3, tagLevel, output) // static|field type varName

	for true {
		words := strings.Split(lines[index], " ")[1]
		word := words[2] // word between the tags
		if word == ',' {
			index = writeLinesToXML(lines, index, 2, tagLevel, output) // (, varName)*
		} else {
			break
		}
	}

	index = writeLinesToXML(lines, index, 1, tagLevel, output) // ;
	tagLevel = tagLevel - 1
	writeToXML("/classVarDec", tagLevel, output)
	return index
}

func writeStatements(lines []string, index int, tagLevel int, output *os.File) int {
	writeToXML("statements", tagLevel, output)
	tagLevel = tagLevel + 1

	for true {
		words := strings.Split(lines[index], " ")[1]
		word := string(words[2]) // word between the tags
		if word == "let" || word == "if" || word == "while" || word == "do" || word == "return" {
			index = writeStatement(lines, index, tagLevel, output) // statement*
		} else {
			break
		}
	}

	tagLevel = tagLevel - 1
	writeToXML("/statements", tagLevel, output)

	return index
}
func writeStatement(lines []string, index int, tagLevel int, output *os.File) int {
	words := strings.Split(lines[index], " ")[1]
	word := words[2] // word between the tags

	switch string(word) { // switch with the word, and goes to right function
	case "let":
		index = writeLetStatement(lines, index, tagLevel, output) // letStatement

	case "if":
		index = writeIfStatement(lines, index, tagLevel, output) // ifStatement

	case "while":
		index = writeWhileStatement(lines, index, tagLevel, output) //whileStatement

	case "do":
		index = writeDoStatement(lines, index, tagLevel, output) // doStatement

	case "return":
		index = writeReturnStatement(lines, index, tagLevel, output) //returnStatement

	default:
		print("No statement here...")

	}

	return index
}
func writeReturnStatement(lines []string, index int, tagLevel int, output *os.File) int {
	writeToXML("returnStatement", tagLevel, output)
	tagLevel = tagLevel + 1

	index = writeLinesToXML(lines, index, 1, tagLevel, output) // return

	words := strings.Split(lines[index], " ")[1]
	word := words[2] // word between the tags
	if word != ';' {
		index = writeExpression(lines, index, tagLevel, output) // expression?
	}

	index = writeLinesToXML(lines, index, 1, tagLevel, output) // ;

	tagLevel = tagLevel - 1
	writeToXML("/returnStatement", tagLevel, output)

	return index
}
func writeDoStatement(lines []string, index int, tagLevel int, output *os.File) int {
	writeToXML("doStatement", tagLevel, output)
	tagLevel = tagLevel + 1

	index = writeLinesToXML(lines, index, 1, tagLevel, output)  // do
	index = writeSubroutineCall(lines, index, tagLevel, output) // subroutineCall
	index = writeLinesToXML(lines, index, 1, tagLevel, output)  // ;

	tagLevel = tagLevel - 1
	writeToXML("/doStatement", tagLevel, output)

	return index
}
func writeWhileStatement(lines []string, index int, tagLevel int, output *os.File) int {
	writeToXML("whileStatement", tagLevel, output)
	tagLevel = tagLevel + 1

	index = writeLinesToXML(lines, index, 2, tagLevel, output) // while (
	index = writeExpression(lines, index, tagLevel, output)    // expression
	index = writeLinesToXML(lines, index, 2, tagLevel, output) // ) {
	index = writeStatements(lines, index, tagLevel, output)    // statements
	index = writeLinesToXML(lines, index, 1, tagLevel, output) // }

	tagLevel = tagLevel - 1
	writeToXML("/whileStatement", tagLevel, output)

	return index
}
func writeIfStatement(lines []string, index int, tagLevel int, output *os.File) int {
	writeToXML("ifStatement", tagLevel, output)
	tagLevel = tagLevel + 1

	index = writeLinesToXML(lines, index, 2, tagLevel, output) // if (
	index = writeExpression(lines, index, tagLevel, output)    // expression
	index = writeLinesToXML(lines, index, 2, tagLevel, output) // ) {
	index = writeStatements(lines, index, tagLevel, output)    // statements
	index = writeLinesToXML(lines, index, 1, tagLevel, output) // }

	words := strings.Split(lines[index], " ")[1]
	word := string(words[2]) // word between the tags
	if word == "else" {      // ?
		index = writeLinesToXML(lines, index, 2, tagLevel, output) // else {
		index = writeStatements(lines, index, tagLevel, output)    // statements
		index = writeLinesToXML(lines, index, 1, tagLevel, output) // }
	}

	tagLevel = tagLevel - 1
	writeToXML("/ifStatement", tagLevel, output)

	return index
}
func writeLetStatement(lines []string, index int, tagLevel int, output *os.File) int {
	writeToXML("letStatement", tagLevel, output)
	tagLevel = tagLevel + 1

	index = writeLinesToXML(lines, index, 2, tagLevel, output) // let varName

	words := strings.Split(lines[index], " ")[1]
	word := words[2] // word between the tags
	if word == '[' { // ?
		index = writeLinesToXML(lines, index, 1, tagLevel, output) // [
		index = writeExpression(lines, index, tagLevel, output)    // expression
		index = writeLinesToXML(lines, index, 1, tagLevel, output) // ]
	}

	index = writeLinesToXML(lines, index, 1, tagLevel, output) // =
	index = writeExpression(lines, index, tagLevel, output)    // expression
	index = writeLinesToXML(lines, index, 1, tagLevel, output) // ;

	tagLevel = tagLevel - 1
	writeToXML("/letStatement", tagLevel, output)

	return index
}
func writeExpression(lines []string, index int, tagLevel int, output *os.File) int {
	writeToXML("expression", tagLevel, output)
	tagLevel = tagLevel + 1

	index = writeTerm(lines, index, tagLevel, output) // term

	for true {
		words := strings.Split(lines[index], " ")[1]
		word := words[2]           // word between the tags
		if TabIsOp(string(word)) { // (op term)*
			index = writeLinesToXML(lines, index, 1, tagLevel, output) // <symbol> op </symbol>
			index = writeTerm(lines, index, tagLevel, output)          // term
		} else {
			break
		}
	}

	tagLevel = tagLevel - 1
	writeToXML("/expression", tagLevel, output)
	return index
}
func writeTerm(lines []string, index int, tagLevel int, output *os.File) int {
	writeToXML("term", tagLevel, output)
	tagLevel = tagLevel + 1

	var words = strings.Split(string(lines[index]), " ")[1]
	var firstWord = words[1] // the opening tag

	switch string(firstWord) { //switch with the word
	case "<integerConstant>":
		index = writeLinesToXML(lines, index, 1, tagLevel, output) // <integerConstant> integer </integerConstant>

	case "<stringConstant>":
		index = writeLinesToXML(lines, index, 1, tagLevel, output) // <stringConstant> string  </stringConstant>

	case "<keyword>":
		var secondWord = string(words[2]) // word between tags
		if TabIsKeyword(secondWord) {
			index = writeLinesToXML(lines, index, 1, tagLevel, output) //<keyword> keywordConstant </keyword>
		} else {
			print("No keyword constant here...")
		}

	case "<identifier>": // multiple options
		var words = strings.Split(string(lines[index+1]), " ")[1] // next line
		var word = words[2]                                       //word between tags

		if word == '[' { // varName [ expression ]
			index = writeLinesToXML(lines, index, 2, tagLevel, output) // varName [
			index = writeExpression(lines, index, tagLevel, output)    // expression
			index = writeLinesToXML(lines, index, 1, tagLevel, output) // ]
		} else if word == '(' || word == '.' { // subroutineCall
			index = writeSubroutineCall(lines, index, tagLevel, output) // subroutineCall
		} else { // varName
			index = writeLinesToXML(lines, index, 1, tagLevel, output) // varName
		}

	case "<symbol>":
		//( or unaryOp
		secondWord := words[2] // word between tags

		if TabIsUnaryOp(string(secondWord)) {
			index = writeLinesToXML(lines, index, 1, tagLevel, output) // <symbol> unaryOp </symbol>
			index = writeTerm(lines, index, tagLevel, output)          //term
		} else if secondWord == '(' {
			index = writeLinesToXML(lines, index, 1, tagLevel, output) //(
			index = writeExpression(lines, index, tagLevel, output)    //expression
			index = writeLinesToXML(lines, index, 1, tagLevel, output) // )
		} else {
			print("No symbol here...")
		}

	default:
		print("No term here...")
	}

	tagLevel = tagLevel - 1
	writeToXML("/term", tagLevel, output)
	return index
}
func writeSubroutineCall(lines []string, index int, tagLevel int, output *os.File) int {
	index = writeLinesToXML(lines, index, 1, tagLevel, output) // identifier

	var words = strings.Split(lines[index], " ")[1] ///////////////////////////////////////////////
	var word = words[2]                             // word between the tags
	if word == '(' {                                //( expressionList )
		index = writeLinesToXML(lines, index, 1, tagLevel, output)  // (
		index = writeExpressionList(lines, index, tagLevel, output) // expressionList
		index = writeLinesToXML(lines, index, 1, tagLevel, output)  // )
	} else if word == '.' { // . identifier ( expressionList )
		index = writeLinesToXML(lines, index, 3, tagLevel, output)  // . identifier (
		index = writeExpressionList(lines, index, tagLevel, output) // expressionList
		index = writeLinesToXML(lines, index, 1, tagLevel, output)  // )
	} else {
		print("ERRADO")
	}
	return index
}
func writeExpressionList(lines []string, index int, tagLevel int, output *os.File) int {
	writeToXML("expressionList", tagLevel, output)
	tagLevel = tagLevel + 1

	var words = strings.Split(string(lines[index]), " ")[1]
	var word = words[2] // word between the tags

	if word != ')' { // not empty list
		index = writeExpression(lines, index, tagLevel, output) // expression

		for true {
			words = strings.Split(string(lines[index]), " ")[1]
			word = words[2]  // word between the tags
			if word == ',' { // (, expression)*
				index = writeLinesToXML(lines, index, 1, tagLevel, output) // ,
				index = writeExpression(lines, index, tagLevel, output)    // expression
			} else {
				break
			}
		}
	}

	tagLevel = tagLevel - 1
	writeToXML("/expressionList", tagLevel, output)

	return index
}

func writeLinesToXML(lines []string, index int, numberOfLines int, levelXML int, output *os.File) int {
	for numberOfLines > 0 {
		writeToXML(lines[index], levelXML, output) //////////////////////
		numberOfLines = numberOfLines - 1
		index = index + 1
	}
	return index
}
func writeToXML(content string, levelXML int, output *os.File) {
	spaces := ""
	for levelXML > 0 {
		spaces += "\t" // tab spaces
		levelXML = levelXML - 1
	}

	if content[0] == '<' {
		var s string = spaces + content
		if _, err := output.WriteString(s); err != nil {
			panic(err)
		}
	} else {
		var s string = spaces + "<" + content + ">"
		if _, err := output.WriteString(s); err != nil {
			panic(err)
		}
	}
}

func TabIsOp(a string) bool {
	var l []string
	l = append(l, "+", "-", "*", "/", "&amp;", "|", "&lt;", "&gt;", "=")

	for _, b := range l {
		if b == a {
			return true
		}
	}
	return false
}
func TabIsUnaryOp(a string) bool {
	var l []byte
	l = append(l, '-', '~')

	for _, b := range l {
		if string(b) == a {
			return true
		}
	}
	return false
}
func TabIsKeyword(a string) bool {
	var l []string
	l = append(l, "true", "false", "null", "this")

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

func Close(a Analyzer) {
	a.InputFile.Close()
	a.OutputFile.Close()
}
