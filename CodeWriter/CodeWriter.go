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

func WriteArithmetic(command string, c CodeWriter, lineNumber string) {
	var splitt []string = strings.Split(command, " ")
	var action string = splitt[0]
	switch action {
	case "add":
		AddTranslate(c)
	case "sub":
		SubTranslate(c)
	case "neg":
		NegTranslate(c)
	case "lt":
		LtTranslate(c, lineNumber)
	case "gt":
		GtTranslate(c, lineNumber)
	case "eq":
		EqTranslate(c, lineNumber)
	case "and":
		AndTranslate(c)
	case "or":
		OrTranslate(c)
	case "not":
		NotTranslate(c)
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
func SubTranslate(c CodeWriter) {
	var s string = "@SP" + "\n" + "M=M-1" + "\n" + "A=M" + "\n" + "D=M" + "\n" + "A=A-1" + "\n" + "M=M-D" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}

}
func EqTranslate(c CodeWriter, lineNumber string) {
	var s string = "@SP" + "\n" + "M=M-1" + "\n" + "A=M" + "\n" + "D=M" + "\n" + "A=A-1" + "\n" + "A=M" + "\n" + "D=A-D" + "\n" +
		"@TRUE" + lineNumber + "\n" + "D;JEQ" + "\n" + "@SP" + "\n" + "A=M-1" + "\n" + "M=0" + "\n" +
		"@FALSE" + lineNumber + "\n" + "0;JEQ" + "\n" +
		"(TRUE" + lineNumber + ")" + "\n" + "@SP" + "\n" + "A=M-1" + "\n" + "M=-1" + "\n" +
		"(FALSE" + lineNumber + ")" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}

}
func NegTranslate(c CodeWriter) {
	var s string = "@SP" + "\n" + "A=M" + "\n" + "A=A-1" + "\n" + "M=-M" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}

}
func LtTranslate(c CodeWriter, lineNumber string) {

	var s string = "@SP" + "\n" + "M=M-1" + "\n" + "A=M" + "\n" + "D=M" + "\n" + "A=A-1" + "\n" +
		"A=M" + "\n" + "D=A-D" + "\n" + "@TRUE" + lineNumber + "\n" + "D;JLT" +
		"\n" + "@SP" + "\n" + "A=M-1" + "\n" + "M=0" + "\n" + "@FALSE" + lineNumber +
		"\n" + "0;JEQ" + "\n" + "(TRUE" + lineNumber + ")" + "\n" + "@SP" + "\n" + "A=M-1" + "\n" + "M=-1" +
		"\n" + "(FALSE" + lineNumber + ")" + "\n" + "\n"

	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}

}
func GtTranslate(c CodeWriter, lineNumber string) {
	var s string = "@SP" + "\n" + "M=M-1" + "\n" + "A=M" + "\n" + "D=M" + "\n" + "A=A-1" + "\n" + "A=M" + "\n" +
		"D=A-D" + "\n" +
		"@TRUE" + lineNumber + "\n" + "D;JGT" + "\n" + "@SP" + "\n" + "A=M-1" + "\n" + "M=0" + "\n" +
		"@FALSE" + lineNumber + "\n" + "0;JEQ" + "\n" +
		"(TRUE" + lineNumber + ")" + "\n" + "@SP" + "\n" + "A=M-1" + "\n" + "M=-1" + "\n" +
		"(FALSE" + lineNumber + ")" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}

}
func AndTranslate(c CodeWriter) {
	var s string = "@SP" + "\n" + "A=M" + "\n" + "A=A-1" + "\n" + "D=M" + "\n" + "A=A-1" + "\n" + "M=M&D" +
		"\n" + "@SP" + "\n" + "M=M-1" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}
func NotTranslate(c CodeWriter) {
	var s string = "@SP" + "\n" + "A=M" + "\n" + "A=A-1" + "\n" + "M=!M" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}

}
func OrTranslate(c CodeWriter) {
	var s string = "@SP" + "\n" + "A=M" + "\n" + "A=A-1" + "\n" + "D=M" + "\n" + "A=A-1" + "\n" +
		"M=M&D" + "\n" + "@SP" + "\n" + "M=M-1" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}

}
func Close(c CodeWriter) {
	c.file.Close()
}
func WritePushPop(command string, segment string, index int, c CodeWriter) {

	if command == "C_PUSH" {
		switch segment {
		case "constant":
			PushConstant(index, c)
		case "local":
			PushLocal(index, c)
		case "temp":
			PushTemp(index, c)
		case "this":
			PushThis(index, c)
		case "that":
			PushThat(index, c)
		case "static":
			PushStatic(index, c)
		case "pointer":
			PushPointer(index, c)
		case "argument":
			PushArgument(index, c)
		}

	}

	if command == "C_POP" {
		switch segment {

		case "local":
			PopLocal(index, c)
		case "temp":
			PopTemp(index, c)
		case "this":
			PopThis(index, c)
		case "that":
			PopThat(index, c)
		case "static":
			PopStatic(index, c)
		case "pointer":
			PopPointer(index, c)
		case "argument":
			PopArgument(index, c)
		}

	}
}

func PushConstant(index int, c CodeWriter) {
	var s string = "@" + strconv.Itoa(index) + "\n" + "D=A" + "\n" + "@SP" + "\n" + "M=M+1" + "\n" + "A=M-1" +
		"\n" + "M=D" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}
func PushLocal(index int, c CodeWriter) {

	var s string = "@" + strconv.Itoa(index) + "\n" + "D=A" + "\n" + "@LCL" + "\n" + "A=M+D" + "\n" + "D=M" + "\n" +
		"@SP" + "\n" + "A=M" + "\n" + "M=D" + "\n" + "@SP" + "\n" + "M=M+1" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}
func PushArgument(index int, c CodeWriter) {
	var s string = "@" + strconv.Itoa(index) + "\n" + "D=A" + "\n" + "@ARG" + "\n" + "A=M+D" + "\n" + "D=M" +
		"\n" + "@SP" + "\n" + "A=M" + "\n" + "M=D" + "\n" + "@SP" + "\n" + "M=M+1" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}
func PushThis(index int, c CodeWriter) {

	var s string = "@" + strconv.Itoa(index) + "\n" + "D=A" + "\n" + "@THIS" + "\n" + "A=M+D" + "\n" + "D=M" +
		"\n" + "@SP" + "\n" + "A=M" + "\n" + "M=D" + "\n" + "@SP" + "\n" + "M=M+1" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}
func PushThat(index int, c CodeWriter) {

	var s string = "@" + strconv.Itoa(index) + "\n" + "D=A" + "\n" + "@THAT" + "\n" + "A=M+D" + "\n" + "D=M" + "\n" +
		"@SP" + "\n" + "A=M" + "\n" + "M=D" + "\n" + "@SP" + "\n" + "M=M+1" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}
func PushTemp(index int, c CodeWriter) {
	////////////////////////////
	var s string = "@" + strconv.Itoa(index+5) + "\n" + "D=M" + "\n" + "@SP" + "\n" + "M=M+1" + "\n" + "A=M-1" +
		"\n" + "M=D" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}
func PushStatic(index int, c CodeWriter) {
	var s string = "@" + c.file.Name() + "." + strconv.Itoa(index) + "\n" + "D=M" + "\n" + "@SP" + "\n" + "A=M" +
		"\n" + "M=D" + "\n" + "@SP" + "\n" + "M=M+1" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}

func PushPointer(index int, c CodeWriter) {
	if index == 0 {
		var s string = "@THIS" + "\n" + "D=M" + "\n" + "@SP" + "\n" + "M=M+1" + "\n" + "A=M-1" + "\n" + "M=D" +
			"\n" + "\n"
		if _, err := c.file.WriteString(s); err != nil {
			panic(err)
		}
	}
	if index == 1 {

		var s string = "@THAT" + "\n" + "D=M" + "\n" + "@SP" + "\n" + "M=M+1" + "\n" + "A=M-1" + "\n" +
			"M=D" + "\n" + "\n"
		if _, err := c.file.WriteString(s); err != nil {
			panic(err)
		}
	}

}
func PopLocal(index int, c CodeWriter) {
	var s string = "@LCL" + "\n" + "D=M" + "\n" + "@" + strconv.Itoa(index) +
		"\n" + "D=D+A" + "\n" + "@R13" + "\n" + "M=D" + "\n" + "@SP" + "\n" +
		"M=M-1" + "\n" + "A=M" + "\n" + "D=M" + "\n" + "@R13" + "\n" + "A=M" +
		"\n" + "M=D" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}

func PopArgument(index int, c CodeWriter) {
	var s string = "@ARG" + "\n" + "D=M" + "\n" + "@" + strconv.Itoa(index) +
		"\n" + "D=D+A" + "\n" + "@R13" + "\n" + "M=D" + "\n" + "@SP" + "\n" +
		"M=M-1" + "\n" + "A=M" + "\n" + "D=M" + "\n" + "@R13" + "\n" + "A=M" +
		"\n" + "M=D" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}

func PopThis(index int, c CodeWriter) {
	var s string = "@THIS" + "\n" + "D=M" + "\n" + "@" + strconv.Itoa(index) +
		"\n" + "D=D+A" + "\n" + "@R13" + "\n" + "M=D" + "\n" + "@SP" + "\n" +
		"M=M-1" + "\n" + "A=M" + "\n" + "D=M" + "\n" + "@R13" + "\n" + "A=M" +
		"\n" + "M=D" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}

func PopThat(index int, c CodeWriter) {
	var s string = "@THAT" + "\n" + "D=M" + "\n" + "@" + strconv.Itoa(index) +
		"\n" + "D=D+A" + "\n" + "@R13" + "\n" + "M=D" + "\n" + "@SP" + "\n" +
		"M=M-1" + "\n" + "A=M" + "\n" + "D=M" + "\n" + "@R13" + "\n" + "A=M" +
		"\n" + "M=D" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}

func PopTemp(index int, c CodeWriter) {
	var s string = "@SP" + "\n" + "A=M-1" + "\n" + "D=M" + "\n" + "@" +
		strconv.Itoa(index+5) + "\n" + "M=D" + "\n" + "@SP" + "\n" + "M=M-1" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}

func PopStatic(index int, c CodeWriter) {
	var s string = "@SP" + "\n" + "A=M-1" + "\n" + "D=M" + "\n" + "@" +
		c.file.Name() + "." + strconv.Itoa(index) + "\n" + "M=D" + "\n" + "@SP" + "\n" + "M=M-1" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}

func PopPointer(index int, c CodeWriter) {
	if index == 0 {
		var s string = "@SP" + "\n" + "M=M-1" + "\n" + "A=M" + "\n" + "D=M" + "\n" + "@THIS" + "\n" + "M=D" + "\n" + "\n"
		if _, err := c.file.WriteString(s); err != nil {
			panic(err)
		}
	} else if index == 1 {
		var s string = "@SP" + "\n" + "M=M-1" + "\n" + "A=M" + "\n" + "D=M" + "\n" + "@THAT" + "\n" + "M=D" + "\n" + "\n"
		if _, err := c.file.WriteString(s); err != nil {
			panic(err)
		}
	}
}
