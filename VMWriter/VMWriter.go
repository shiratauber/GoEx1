package VMWriter

import (
	"log"
	"os"
	"strconv"
)

type Writer struct {
	OutputFile *os.File
}

func New(outputFile *os.File) Writer {

	//create the output file
	//outputFile, err := os.Create(output + ".vm")
	//Check(err)
	//open the output file
	//myFile, err := os.OpenFile(outputFile.Name(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	//Check(err)

	mofa := Writer{outputFile}
	return mofa

}

func WritePush(segment string, index int, c Writer) {
	var s string = "push " + segment + " " + strconv.Itoa(index) + "\n"
	if _, err := c.OutputFile.WriteString(s); err != nil {
		panic(err)
	}

}

func WritePop(segment string, index int, c Writer) {
	var s string = "pop " + segment + " " + strconv.Itoa(index) + "\n"
	if _, err := c.OutputFile.WriteString(s); err != nil {
		panic(err)
	}

}

func WriteCommand(command string, c Writer) {
	var s string = command + "\n"
	if _, err := c.OutputFile.WriteString(s); err != nil {
		panic(err)
	}

}

func WriteLabel(label string, c Writer) {
	var s string = label + "\n"
	if _, err := c.OutputFile.WriteString(s); err != nil {
		panic(err)
	}

}

// Writes a VM arithmetic command.
func WriteArithmetic(command string, c Writer) {
	var s string = command + "\n"
	if _, err := c.OutputFile.WriteString(s); err != nil {
		panic(err)
	}
}

func WriteGoto(label string, c Writer) {
	var s string = "goto " + label + "\n"
	if _, err := c.OutputFile.WriteString(s); err != nil {
		panic(err)
	}

}

func WriteIf(label string, c Writer) {
	var s string = "if-goto " + label + "\n"
	if _, err := c.OutputFile.WriteString(s); err != nil {
		panic(err)
	}

}

func WriteCall(name string, nArgs int, c Writer) {
	var s string = "call " + name + " " + strconv.Itoa(nArgs) + "\n"
	if _, err := c.OutputFile.WriteString(s); err != nil {
		panic(err)
	}

}

func WriteFunction(name string, nLocals int, c Writer) {
	var s string = "function " + name + " " + strconv.Itoa(nLocals) + "\n"
	if _, err := c.OutputFile.WriteString(s); err != nil {
		panic(err)
	}

}

func WriteReturn(c Writer) {
	var s string = "return" + "\n"
	if _, err := c.OutputFile.WriteString(s); err != nil {
		panic(err)
	}

}

func Close(c Writer) {

	c.OutputFile.Close()
}
func Check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
