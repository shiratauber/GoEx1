package CodeWriter

import (
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

type CodeWriter struct {
	outputFile  string
	file        *os.File
	current     string
	CallCounter int
	funcName    string
}

func New(path string) CodeWriter {
	var splitt []string = strings.Split(path, "\\")
	var last string = splitt[len(splitt)-1] //the name of the output file

	//create the output file
	oFile, err := os.Create(last + ".asm")
	Check(err)
	//open the output file
	myFile, err := os.OpenFile(oFile.Name(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	Check(err)
	output := CodeWriter{path, myFile, " ", 0, last}
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
		"M=M|D" + "\n" + "@SP" + "\n" + "M=M-1" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}

}
func Close(c CodeWriter) {
	c.file.Close()
}
func WritePushPop(command string, segment string, index int, c CodeWriter, staticScope string) {

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
			PushStatic(index, c, staticScope)
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
			PopStatic(index, c, staticScope)
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
func PushStatic(index int, c CodeWriter, staticScope string) {
	var s string = "@" + staticScope + "." + strconv.Itoa(index) + "\n" + "D=M" + "\n" + "@SP" + "\n" + "A=M" +
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
	//var s string = "@SP" + "\n" + "A=M-1" + "\n" + "D=M" + "\n" + "@" +
	//strconv.Itoa(index+5) + "\n" + "M=D" + "\n" + "@SP" + "\n" + "M=M-1" + "\n" + "\n"
	var s string = "@5" + "\n" + "D=A" + "\n" + "@" + strconv.Itoa(index) + "\n" + "D=D+A" + "\n" + "@13" + "\n" + "M=D" + "\n" + "@SP" + "\n" + "M=M-1" + "\n" +
		"A=M" + "\n" + "D=M" + "\n" + "@13" + "\n" + "A=M" + "\n" + "M=D" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}

func PopStatic(index int, c CodeWriter, staticScope string) {
	var s string = "@SP" + "\n" + "A=M-1" + "\n" + "D=M" + "\n" + "@" +
		staticScope + "." + strconv.Itoa(index) + "\n" + "M=D" + "\n" + "@SP" + "\n" + "M=M-1" +
		"\n" + "\n"
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

//EX2
//////////////////////////////////////////////////////////
func WriteInit(haveSysInit bool, c *CodeWriter) {
	var s string = "@256" + "\n" + "D=A" + "\n" + "@SP" + "\n" + "M=D" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
	if haveSysInit == true {
		c.CallCounter = 0
		WriteCall("Sys.init", "0", *c)
	}
}
func WriteLabel(arg1 string, c CodeWriter) {
	var s string = "(" + c.file.Name() + "." + arg1 + ")" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}

func WriteGoTo(arg1 string, c CodeWriter) {
	var s string = "@" + c.file.Name() + "." + arg1 + "\n" + "0;JMP" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}

func WriteIf(arg1 string, c CodeWriter) {
	var s string = "@SP" + "\n" + "M=M-1" + "\n" + "A=M" + "\n" + "D=M" + "\n" + "@" + c.file.Name() + "." + arg1 +
		"\n" + "D;JNE" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}

func WriteFunction(arg1 string, arg2 string, c CodeWriter) {
	c.funcName = arg1
	//label f
	var s string = "(" + arg1 + ")" + "\n"
	//WriteLabel(arg1, c)
	//initialize local variables
	s += "@" + arg2 + "\n" + "D=A" + "\n" + "@" + arg1 + ".END" + "\n" + "D;JEQ" + "\n"
	//jump if false- k!=0
	s += "(" + arg1 + ".LOOP)" + "\n" + "@SP" + "\n" + "A=M" + "\n" + "M=0" + "\n" + "@SP" + "\n" + "M=M+1" +
		"\n" + "@" + arg1 + ".LOOP" + "\n"
	//jump while k!=0
	s += "D=D-1;JNE" + "\n"
	//jump if true k==0
	s += "(" + arg1 + ".END)" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}

func WriteCall(arg1 string, arg2 string, c CodeWriter) {

	// push return-address
	var s string = "@" + arg1 + ".RETURN_ADDRESS" + strconv.Itoa(c.CallCounter) + "\n" + "D=A" + "\n" + "@SP" + "\n" +
		"A=M" + "\n" + "M=D" + "\n" + "@SP" + "\n" + "M=M+1" + "\n"
	// push LCL
	s += "@LCL" + "\n" + "D=M" + "\n" + "@SP" + "\n" + "A=M" + "\n" + "M=D" + "\n" + "@SP" + "\n" + "M=M+1" + "\n"
	// push ARG
	s += "@ARG" + "\n" + "D=M" + "\n" + "@SP" + "\n" + "A=M" + "\n" + "M=D" + "\n" + "@SP" + "\n" + "M=M+1" + "\n"
	// push THIS
	s += "@THIS" + "\n" + "D=M" + "\n" + "@SP" + "\n" + "A=M" + "\n" + "M=D" + "\n" + "@SP" + "\n" + "M=M+1" + "\n"
	// push THAT
	s += "@THAT" + "\n" + "D=M" + "\n" + "@SP" + "\n" + "A=M" + "\n" + "M=D" + "\n" + "@SP" + "\n" + "M=M+1" + "\n"
	// ARG = SP-n-5
	y, e := strconv.Atoi(arg2)
	Check(e)
	s += "@SP" + "\n" + "D=M" + "\n" + "@" + strconv.Itoa(y+5) + "\n" + "D=D-A" + "\n" + "@ARG" + "\n" + "M=D" + "\n"
	// LCL = SP
	s += "@SP" + "\n" + "D=M" + "\n" + "@LCL" + "\n" + "M=D" + "\n"
	// goto f
	s += "@" + arg1 + "\n" + "0;JMP" + "\n"
	// label return-address
	s += "(" + arg1 + ".RETURN_ADDRESS" + strconv.Itoa(c.CallCounter) + ")" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
	//c.callCounter++
}
func WriteReturn(c CodeWriter) {
	// FRAME = LCL
	var s string = "@LCL" + "\n" + "D=M" + "\n"
	// RET = *(FRAME - 5)
	// RAM[13] = (LOCAL - 5)
	s += "@5" + "\n" + "A=D-A" + "\n" + "D=M" + "\n" + "@13" + "\n" + "M=D" + "\n"
	// *ARG = pop()
	s += "@SP" + "\n" + "M=M-1" + "\n" + "A=M" + "\n" + "D=M" + "\n" + "@ARG" + "\n" + "A=M" + "\n" + "M=D" + "\n"
	// SP = ARG + 1
	s += "@ARG" + "\n" + "D=M" + "\n" + "@SP" + "\n" + "M=D+1" + "\n"
	// THAT = *(FRAME - 1)
	s += "@LCL" + "\n" + "M=M-1" + "\n" + "A=M" + "\n" + "D=M" + "\n" + "@THAT" + "\n" + "M=D" + "\n"
	// THIS = *(FRAME - 2)
	s += "@LCL" + "\n" + "M=M-1" + "\n" + "A=M" + "\n" + "D=M" + "\n" + "@THIS" + "\n" + "M=D" + "\n"
	// ARG = *(FRAME - 3)
	s += "@LCL" + "\n" + "M=M-1" + "\n" + "A=M" + "\n" + "D=M" + "\n" + "@ARG" + "\n" + "M=D" + "\n"
	// LCL = *(FRAME - 4)
	s += "@LCL" + "\n" + "M=M-1" + "\n" + "A=M" + "\n" + "D=M" + "\n" + "@LCL" + "\n" + "M=D" + "\n"
	// goto RET
	s += "@13" + "\n" + "A=M" + "\n" + "0;JMP" + "\n" + "\n"
	if _, err := c.file.WriteString(s); err != nil {
		panic(err)
	}
}

/*
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
	output := CodeWriter{path, myFile, " ", 0}
	return output
}*/
