package CodeWriter

import (
	"log"
	"os"
)

type CodeWriter struct {
	outputFile string
	file       *os.File
	current    string
}

func New(path string) CodeWriter {
	oFile, err := os.Create(path + ".asm")
	check(err)
	//open the output file
	myFile, err := os.OpenFile(oFile.Name(), os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0777)
	check(err)
	output := CodeWriter{path, myFile, " "}
	return output
}

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
