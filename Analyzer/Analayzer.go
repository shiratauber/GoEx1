package Analyzer

import (
	"go/ast"
	"os"
)

type Analyzer struct {
	outputFile  string
	file        *os.File
	current     string
	CallCounter int
	funcName    string
}

func New(path string) Analyzer {
	ast.SelectStmt{Select: ast.SelectStmt
	}
	return Analyzer{}
}

func ss() {
	var stamm string = " "
	println(stamm)
}

func searchFiles(pathToSearch string) {

}
func writeClass (lines string, index int, output *os.File) {

}
func writeSubroutineDec(lines string, index int, tagLevel int, output *os.File) {

}
func writeSubroutineBody(lines string, index int, tagLevel int , output *os.File){

}
func writeVarDec(lines string, index int, tagLevel int, output string){

}
func writeParameterList(lines string, index int, tagLevel int, output string){


}
func writeClassVarDec(lines string, index int, tagLevel int, output string) {

}

func writeStatements(lines string, index int, tagLevel int, output string){

}
func writeStatement(lines string, index int, tagLevel int, output string) {

}
func writeReturnStatement(lines string, index int, tagLevel int, output string) {

}
func writeDoStatement(lines string, index int, tagLevel int, output string) {

}
func writeWhileStatement(lines string, index int, tagLevel int, output string) {

}
func writeIfStatement(lines string, index int, tagLevel int, output string) {

}
func writeLetStatement(lines string, index int, tagLevel int, output string) {

}
func writeExpression(lines string, index int, tagLevel int, output string) {

}
func writeTerm(lines string, index int, tagLevel int, output string) {

}
func writeSubroutineCall(lines string, index int, tagLevel int, output string) {

}
func writeExpressionList(lines string, index int, tagLevel int, output string) {

}

func writeLinesToXML(lines string, index int,numberOfLines int, levelXML int, output string) {

}
func writeToXML(content string, levelXML int, output *os.File){

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
