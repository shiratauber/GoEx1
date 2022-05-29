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

func writePush(segment int, index int, c Writer) {
	var s string = "push " + strconv.Itoa(segment) + " " + strconv.Itoa(index) + "\n"
	if _, err := c.OutputFile.WriteString(s); err != nil {
		panic(err)
	}

}

func writePop(segment int, index int, c Writer) {
	var s string = "pop " + strconv.Itoa(segment) + " " + strconv.Itoa(index) + "\n"
	if _, err := c.OutputFile.WriteString(s); err != nil {
		panic(err)
	}

}

func writeCommand(command string, c Writer) {
	var s string = command + "\n"
	if _, err := c.OutputFile.WriteString(s); err != nil {
		panic(err)
	}

}

func writeLabel(label string, c Writer) {
	var s string = label + "\n"
	if _, err := c.OutputFile.WriteString(s); err != nil {
		panic(err)
	}

}

func writeGoto(label string, c Writer) {
	var s string = "goto " + label + "\n"
	if _, err := c.OutputFile.WriteString(s); err != nil {
		panic(err)
	}

}

func writeIf(label string, c Writer) {
	var s string = "if-goto " + label + "\n"
	if _, err := c.OutputFile.WriteString(s); err != nil {
		panic(err)
	}

}

func writeCall(name string, nArgs int, c Writer) {
	var s string = "call " + name + " " + strconv.Itoa(nArgs) + "\n"
	if _, err := c.OutputFile.WriteString(s); err != nil {
		panic(err)
	}

}

func writeFunction(name string, nLocals int, c Writer) {
	var s string = "function " + name + " " + strconv.Itoa(nLocals) + "\n"
	if _, err := c.OutputFile.WriteString(s); err != nil {
		panic(err)
	}

}

func writeReturn(c Writer) {
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
