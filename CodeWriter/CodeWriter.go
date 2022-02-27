package CodeWriter

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type CodeWriter struct {
	outputFile string
	file       *os.File
	current    string
}

func New(path string) CodeWriter {
	var split1 []string = strings.Split(path, ".")
	var withoutLast []string = split1[:len(split1)-1]
	//var withLast string = strings.Join(withoutLast," ")
	var split2 []string = strings.Split(withoutLast[0], "\\")
	var last string = split2[len(split2)-1]
	oFile, err := os.Create(last + ".asm")
	Check(err)
	//open the output file
	myFile, err := os.OpenFile(oFile.Name(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	Check(err)
	output := CodeWriter{path, myFile, " "}
	return output
}
func SetFileName(c CodeWriter, s string) {
	c.outputFile = s

}
func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func WriteArithmetic(command string, c CodeWriter) {
	var splitt []string = strings.Split(command, " ")
	var action string = splitt[0]
	switch action {
	case "add":
		AddTranslate(c)
	case "sub":
		fmt.Println("Tuesday")
	case "neg":
		fmt.Println("Wednesday")
	case "lt":
		fmt.Println("Thursday")
	case "gt":
		fmt.Println("Friday")
	case "eq":
		fmt.Println("Saturday")
	case "and":
		fmt.Println("Sunday")
	case "or":
		fmt.Println("Sunday")
	case "not":
		fmt.Println("Sunday")
	default:
		fmt.Println("Invalid")
	}
}
func AddTranslate(c CodeWriter) {

	var s string = "@SP" + "\n" + "M=M-1" + "\n" + "A=M" + "\n" + "D=M" + "\n" + "A=A-1" + "\n" + "M=M+D" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}

}
func SubTranslate() {

}
func EqTranslate() {

}
func NegTranslate() {

}
func LtTranslate() {

}
func GtTranslate() {

}
func AndTranslate() {

}
func NotTranslate() {

}
func OrTranslate() {

}
func Close(c CodeWriter) {
	c.file.Close()
}
func WritePushPop(command string, segment string, index int, c CodeWriter) {

	if command == "C_PUSH" {
		switch segment {
		case "constant":
			PushConstant(index, c)

		}
	}

}

func PushConstant(index int, c CodeWriter) {
	var s string = "@" + strconv.Itoa(index) + "\n" + "D=A" + "\n" + "@SP" + "\n" + "M=M+1" + "\n" + "A=M-1" + "\n" + "M=D" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}
