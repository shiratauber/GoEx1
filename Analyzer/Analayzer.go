package Analyzer

import (
	"os"
	"strings"
)

type Analyzer struct {
	outputFile  string
	file        *os.File
	current     string
	CallCounter int
	funcName    string
}

func New(path string) Analyzer {

	return Analyzer{}
}

func ss() {
	var stamm string = " "
	println(stamm)
}

func searchFiles(pathToSearch string) {

}
func writeClass(lines string, index int, output *os.File) {

}
func writeSubroutineDec(lines string, index int, tagLevel int, output *os.File) int {
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
func writeSubroutineBody(lines string, index int, tagLevel int, output *os.File) int {
	writeToXML("subroutineBody", tagLevel, output)
	tagLevel = tagLevel + 1

	index = writeLinesToXML(lines, index, 1, tagLevel, output)  {

	for true{
	var words = strings.Split(string(lines[index]), " ")[1]////////////////////////////////////////////
	var word = words[2]  //word between the tags

	if(word == "var"){
	index = writeVarDec(lines, index, tagLevel, output) // varDec*
	}else{
	break
	}
	}

	index = writeStatements(lines, index, tagLevel, output) // statements
	index = writeLinesToXML(lines, index, 1, tagLevel, output)  }

	tagLevel = tagLevel - 1
	writeToXML("/subroutineBody", tagLevel, output)

	return index
}
func writeVarDec(lines string, index int, tagLevel int, output *os.File) int {

}
func writeParameterList(lines string, index int, tagLevel int, output *os.File) int {
	return 0
}
func writeClassVarDec(lines string, index int, tagLevel int, output *os.File) {

}

func writeStatements(lines string, index int, tagLevel int, output *os.File) int{

}
func writeStatement(lines string, index int, tagLevel int, output *os.File) {

}
func writeReturnStatement(lines string, index int, tagLevel int, output *os.File) {

}
func writeDoStatement(lines string, index int, tagLevel int, output *os.File) {

}
func writeWhileStatement(lines string, index int, tagLevel int, output *os.File) {

}
func writeIfStatement(lines string, index int, tagLevel int, output *os.File) {

}
func writeLetStatement(lines string, index int, tagLevel int, output *os.File) {

}
func writeExpression(lines string, index int, tagLevel int, output *os.File) int {

}
func writeTerm(lines string, index int, tagLevel int, output *os.File) {
	writeToXML("term", tagLevel, output)
	tagLevel = tagLevel + 1

	var words = strings.Split(string(lines[index]), " ")[1]
	var firstWord = words[1]  // the opening tag

	switch firstWord { //switch with the word
	case "<integerConstant>":
		index = writeLinesToXML(lines, index, 1, tagLevel, output) // <integerConstant> integer </integerConstant>

	case "<stringConstant>":
		index = writeLinesToXML(lines, index, 1, tagLevel, output) // <stringConstant> string  </stringConstant>

	case "<keyword>":
		var secondWord = words[2] // word between tags
		if (secondWord in keywordConstant) {
		index = writeLinesToXML(lines, index, 1, tagLevel, output) //<keyword> keywordConstant </keyword>
	} else {
		print("No keyword constant here...")
	}
	}
	case "<identifier>": // multiple options
	words = strsplit(lines[index + 1], " ")[[1]] // next line
	word = words[2] //word between tags

	if (word == '[') { // varName [ expression ]
	index <- writeLinesToXML(lines, index, 2, tagLevel, output) // varName [
	index <- writeExpression(lines, index, tagLevel, output) // expression
	index <- writeLinesToXML(lines, index, 1, tagLevel, output) // ]
	} else if (word == '(' | word == '.') { // subroutineCall
	index <- writeSubroutineCall(lines, index, tagLevel, output) // subroutineCall
	} else { // varName
	index <- writeLinesToXML(lines, index, 1, tagLevel, output) // varName
	}
	},
	"<symbol>"={ //( or unaryOp
	secondWord <- words[2] // word between tags

	if (secondWord %in% unaryOp) {
	index =writeLinesToXML(lines, index, 1, tagLevel, output) // <symbol> unaryOp </symbol>
	index = writeTerm(lines, index, tagLevel, output) //term
	} else if (secondWord == '(') {
	index = writeLinesToXML(lines, index, 1, tagLevel, output) //(
	index = writeExpression(lines, index, tagLevel, output) //expression
	index = writeLinesToXML(lines, index, 1, tagLevel, output) // )
	} else {
	print("No symbol here...")
	}
	},
	 default:
	{
	print("No term here...")
	}
}

	tagLevel = tagLevel - 1
	writeToXML("/term", tagLevel, output)

	return index
}
func writeSubroutineCall(lines string, index int, tagLevel int, output *os.File) int {
	index = writeLinesToXML(lines, index, 1, tagLevel, output) // identifier

	var words = strings.Split(string(lines[index]), " ")[1]///////////////////////////////////////////////
	var word = words[2]  // word between the tags
	if(word == '('){ //( expressionList )
	index = writeLinesToXML(lines, index, 1, tagLevel, output) // (
	index = writeExpressionList(lines, index, tagLevel, output) // expressionList
	index = writeLinesToXML(lines, index, 1, tagLevel, output) // )
	} else if(word == '.'){ // . identifier ( expressionList )
	index = writeLinesToXML(lines, index, 3, tagLevel, output) // . identifier (
	index = writeExpressionList(lines, index, tagLevel, output) // expressionList
	index = writeLinesToXML(lines, index, 1, tagLevel, output) // )
	} else{
	print("ERRADO")
	}

	return index
}
func writeExpressionList(lines string, index int, tagLevel int, output *os.File) int {
	writeToXML("expressionList", tagLevel, output)
	tagLevel = tagLevel + 1

	var words = strings.Split(string(lines[index]), " ")[1]
	var word = words[2]  // word between the tags

	if(word != ')'){   // not empty list, ?
	index = writeExpression(lines, index, tagLevel, output) // expression

	for true{
	words = strings.Split(string(lines[index]), " ")[1]
	word= words[2]  // word between the tags
	if(word == ','){ // (, expression)*
	index = writeLinesToXML(lines, index, 1, tagLevel, output) // ,
	index = writeExpression(lines, index, tagLevel, output) // expression
	} else{
	break
	}
	}
	}

	tagLevel = tagLevel - 1
	writeToXML("/expressionList", tagLevel, output)

	return index
}

func writeLinesToXML(lines string, index int, numberOfLines int, levelXML int, output *os.File) int {
	for true {
		if (numberOfLines > 0) {
			writeToXML(string(lines[index]), levelXML, output) //////////////////////
			numberOfLines = numberOfLines - 1
			index = index + 1
		}else{
			break
		}
	}

	return index
}
func writeToXML(content string, levelXML int, output *os.File) {

}

func TabIsOp(a string) bool {
	var l []byte
	l = append(l, '0', '1', '2', '3', '4', '5', '6', '7', '8', '9')

	for _, b := range l {
		if b == a {
			return true
		}
	}
	return false
}
func TabIsUnaryOp(a byte) bool {
	var l []byte
	l = append(l, '0', '1', '2', '3', '4', '5', '6', '7', '8', '9')

	for _, b := range l {
		if b == a {
			return true
		}
	}
	return false
}
func TabIsKeyword(a string) bool {
	var l []byte
	l = append(l, '0', '1', '2', '3', '4', '5', '6', '7', '8', '9')

	for _, b := range l {
		if b == a {
			return true
		}
	}
	return false
}
